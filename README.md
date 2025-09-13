# Commands


```
go mod init github.com/aayushxrj/gRPC-streaming-demo
go get google.golang.org/grpc
go mod tidy
```
```
protoc --go_out=. --go-grpc_out=. proto/main.proto
```