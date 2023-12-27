# `smbsocks`/`rpclink`

```
a simple Go-only rpc2socks alternative. the client is windows-only 
and should be run with runas /netonly with administrative credentials

establishes SOCKS over SSH over named pipes.
vendors npipe to make the pipe remotely accessible.
```
```
usage:

mingw32-make

cmd > rundll32 rpcsrv.dll,DllGetClassObject
# then use a windows extension in Sliver and portforward the local port
sliver > rpclink \\\\<ip>\\pipe\\port
sliver > portfwd add -b 127.0.0.1:9050 -r 127.0.0.1:9050
```
