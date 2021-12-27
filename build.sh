export CGO_ENABLED=0
GOOS=linux GOARCH=amd64 go build -o lmp-linux64
GOOS=darwin GOARCH=amd64 go build -o lmp-darwin64
GOOS=windows GOARCH=amd64 go build -o lmp-windows64.exe
GOOS=windows GOARCH=386 go build -o lmp-windows32.exe
GOOS=freebsd GOARCH=amd64 go build -o lmp-freebsd64
GOOS=openbsd GOARCH=amd64 go build -o lmp-openbsd64
shasum -a 256 lmp-* > checksums.txt