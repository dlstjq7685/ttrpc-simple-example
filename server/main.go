package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	example "ttrpc/simple/example/server/interfaces"

	"github.com/containerd/ttrpc"
)

const socket = "50051"

func main() {
	log.Println("server start")
	if err := handle(); err != nil {
		log.Fatal(err)
	}
}

func handle() error {
	return server()
}

func serverIntercept(ctx context.Context, um ttrpc.Unmarshaler, i *ttrpc.UnaryServerInfo, m ttrpc.Method) (interface{}, error) {
	log.Println("server interceptor")
	dumpMetadata(ctx)
	return m(ctx, um)
}

func dumpMetadata(ctx context.Context) {
	md, ok := ttrpc.GetMetadata(ctx)
	if !ok {
		panic("no metadata")
	}
	fmt.Println(md)
	if err := json.NewEncoder(os.Stdout).Encode(md); err != nil {
		panic(err)
	}
}

func server() error {
	s, err := ttrpc.NewServer(
		ttrpc.WithServerHandshaker(nil),
		ttrpc.WithUnaryServerInterceptor(serverIntercept),
	)
	if err != nil {
		return err
	}
	defer s.Close()
	example.RegisterExampleService(s, &exampleServer{})

	l, err := net.Listen("unix", socket)
	if err != nil {
		return err
	}
        defer func() {
		l.Close()
		os.Remove(socket)
	}()

	return s.Serve(context.Background(), l)
}

type exampleServer struct {
}

func (s *exampleServer) Method1(ctx context.Context, r *example.Method1Request) (*example.Method1Response, error) {
	fmt.Println("method1 call")
	return &example.Method1Response{
		Foo: r.Foo,
		Bar: r.Bar,
	}, nil
}
