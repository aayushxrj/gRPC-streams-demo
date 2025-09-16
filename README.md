# Commands


```
go mod init github.com/aayushxrj/gRPC-streaming-demo
go get google.golang.org/grpc
go mod tidy
```
```
protoc --go_out=. --go-grpc_out=. proto/main.proto
```

# protoc-gen-validate

```
go install github.com/envoyproxy/protoc-gen-validate@latest
```
```
nano .bash_rc
```
```
export PATH="$PATH:$(go env GOPATH)/bin"
```
```
source ~/.bashrc
```

Include validate/validate.proto or do this
```
mkdir -p proto/validate
curl -sSL https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/main/validate/validate.proto -o proto/validate/validate.proto
```

```
protoc \
  -I proto \
  --go_out=proto/gen --go_opt=paths=source_relative \
  --go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative \
  --validate_out="lang=go,paths=source_relative:proto/gen" \
  proto/main.proto
```

```
go get github.com/envoyproxy/protoc-gen-validate/validate
```
