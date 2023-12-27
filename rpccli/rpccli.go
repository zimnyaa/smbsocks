package main

import (
    "golang.org/x/crypto/ssh"
    "net"
    "fmt"
    "log"
    "context"
    "github.com/armon/go-socks5"
    "github.com/zimnyaa/smbsocks/npipe"
)
import "C"

func findUnusedPort(startPort int32) (int32) {
    for port := startPort; port <= 65535; port++ {
        addr := fmt.Sprintf("localhost:%d", port)
        listener, err := net.Listen("tcp", addr)
        if err != nil {
            continue
        }
        listener.Close()
        return port
    }
    return 0
}

type sshResolver struct{
    sshConnection *ssh.Client
}


func (d sshResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {

    sess, err := d.sshConnection.NewSession()
    if err != nil {
        return ctx, nil, fmt.Errorf("sess err.")
    }
    defer sess.Close()
    stdin, err := sess.StdinPipe()
    if err != nil {
        return ctx, nil, fmt.Errorf("pipe err.")
    }

    stdout, err := sess.StdoutPipe()
    if err != nil {
        return ctx, nil, fmt.Errorf("pipe err.")
    }

    stdin.Write([]byte(name))
    defer stdin.Close()
    var addr []byte = make([]byte, 256) 
    
    _, err = stdout.Read(addr)
    if err != nil {
        return ctx, nil, fmt.Errorf("pipe err.")
    }

    resp := string(addr)

    if resp == "err" {
        return ctx, nil, fmt.Errorf("resolve err.")
    }
    ipaddr := net.ParseIP(resp)
    return ctx, ipaddr, err
}

//export dialpipe
func dialpipe(urlc *C.char) {
    fmt.Println("go received args:", urlc)
    url := C.GoString(urlc)
    fmt.Println("bind url:", url)

    socksconn, err := npipe.Dial(url)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer socksconn.Close()

    sshConf := &ssh.ClientConfig{
        User:            "root",
        Auth:            []ssh.AuthMethod{ssh.Password("asdf")},
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }
    
    c, chans, reqs, err := ssh.NewClientConn(socksconn, "255.255.255.255", sshConf)
    if err != nil {
        log.Printf("%v", err)
        return
    }
    sshConn := ssh.NewClient(c, chans, reqs)
    sshRes := sshResolver{sshConnection: sshConn}
    
    defer sshConn.Close()

    //log.Printf("connected to backwards ssh server\n")

    conf := &socks5.Config{
        Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
            return sshConn.Dial(network, addr)
        },
        Resolver: sshRes,
    }

    serverSocks, err := socks5.New(conf)
    if err != nil {
        fmt.Println(err)
        return
    }
    port := findUnusedPort(9050)
    //log.Printf("creating a socks server@%d\n", port)
    if err := serverSocks.ListenAndServe("tcp", fmt.Sprintf("0.0.0.0:%d", port)); err != nil {
        log.Fatalf("failed to create socks5 server%v\n", err)
    }

    return
}

func main() {}