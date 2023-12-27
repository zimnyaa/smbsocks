.PHONY: all

all:
	GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 go build -ldflags "-s -w" --buildmode=c-shared -o rpcsrv.dll ./rpcsrv/
	GOOS=windows GOARCH=amd64 go build -o rpccli.exe ./rpccli/
