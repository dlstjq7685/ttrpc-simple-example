# ttrpc-simple-example


## 1. prepare ttrpc plugin

`go install github.com/containerd/ttrpc/cmd/protoc-gen-go-ttrpc@v1.2.3`
`go install github.com/containerd/ttrpc/cmd/protoc-gen-gogottrpc@v1.2.3`


## 2. generate ttrpc script

`protoc --go_out=. --go_opt=paths=source_relative --go-ttrpc_out=. --go-ttrpc_opt=paths=source_relative proto/example.proto`

## 3. server, client execute
