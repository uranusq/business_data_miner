# Not working
set GOARCH=amd64
set GOOS=linux
set CGO_ENABLED=1
set CC=x86_64-w64-mingw32-gcc
set CXX=x86_64-w64-mingw32-c++
go build -ldflags "-linkmode external -extldflags -static"