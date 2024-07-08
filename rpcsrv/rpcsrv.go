package main

import (
	"golang.org/x/crypto/ssh"
	"net"
	"io"
	"log"
	"crypto/rsa"
	"crypto/rand"
	"fmt"
	"github.com/zimnyaa/smbsocks/npipe"
)

import "C"

//export DllGetClassObject
func DllGetClassObject() {
	pipename := "testpipename"

	l, err := npipe.Listen(fmt.Sprintf("\\\\.\\pipe\\%s", pipename))
	if err != nil {
		log.Fatalf("listen err: %v", err)
	}

	for {

		nConn, err := l.Accept()
		log.Printf("new conn %v\n", nConn.RemoteAddr)
			
		config := &ssh.ServerConfig{
			PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
				return nil, nil
			},
		}
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			panic(err)
		}
		signer, err := ssh.NewSignerFromKey(privateKey)
		if err != nil {
			panic("Failed to create signer")
		}
		config.AddHostKey(signer)

		sshConn, chans, reqs, err := ssh.NewServerConn(nConn, config)
		if err != nil {
			log.Fatalf("Failed to open stream: %v", err)
		}
		defer sshConn.Close()


		go ssh.DiscardRequests(reqs)

		for newChannel := range chans {
			log.Printf("[socks] new channel\n")
			if newChannel.ChannelType() == "session" {
				go func() {
					connection, requests, err := newChannel.Accept()
					if err != nil {
						return
					}
					go ssh.DiscardRequests(requests)
					var domainBytes []byte = make([]byte, 1024)
					n, err := connection.Read(domainBytes)
					if err != nil || n == 0 {
						return
					}
					connection.Write(dnsResolve(string(domainBytes)))
					connection.Close()
				}()
				continue
			}

			if newChannel.ChannelType() != "direct-tcpip" {
				newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
				continue
			}
			

			var dReq struct {
				DestAddr string
				DestPort uint32
			}
			ssh.Unmarshal(newChannel.ExtraData(), &dReq)

			log.Printf("new direct-tcpip channel to %s:%d\n", dReq.DestAddr, dReq.DestPort)
			go func() {
				dest := fmt.Sprintf("%s:%d", dReq.DestAddr, dReq.DestPort)
				var conn net.Conn
				var err error
				conn, err = net.Dial("tcp", dest)
					
				if err == nil {
					channel, chreqs, err := newChannel.Accept()
					if err != nil {
						return
					}
					go ssh.DiscardRequests(chreqs)
		
					go func() {
						defer channel.Close()
						defer conn.Close()
						io.Copy(channel, conn)
					}()
					go func() {
						defer channel.Close()
						defer conn.Close()
						io.Copy(conn, channel)
					}()
				}
			}()
		}
	}
	
}

func dnsResolve(name string) ([]byte) {
	log.Printf("dnsresolve: %s\n", name)
	addr, err := net.ResolveIPAddr("ip", name)
	if err != nil {
		return []byte("err")
	}
	return []byte(addr.IP.String())
}

func main() {}
