package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	example "ttrpc/simple/example/client/interfaces"

	"github.com/containerd/ttrpc"
)

const socket = "../server/50051"

func main() {
	log.Println("client start")
	if err := handle(); err != nil {
		log.Fatal(err)
	}
}

func handle() error {
	return client()

}

func clientIntercept(ctx context.Context, req *ttrpc.Request, resp *ttrpc.Response, i *ttrpc.UnaryClientInfo, invoker ttrpc.Invoker) error {
	log.Println("client interceptor")
	dumpMetadata(ctx)
	return invoker(ctx, req, resp)
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
func client() error {
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return err
	}
	defer conn.Close()

	tc := ttrpc.NewClient(conn, ttrpc.WithUnaryClientInterceptor(clientIntercept))
	client := example.NewExampleClient(tc)

	r := &example.Method1Request{
		Foo: "Foo",
		Bar: "Bar",
	}

	ctx := context.Background()
	md := ttrpc.MD{}
	md.Set("name", "koye")
	ctx = ttrpc.WithMetadata(ctx, md)

	resp, err := client.Method1(ctx, r)
	if err != nil {
		log.Println("error")
		return err
	}
	return json.NewEncoder(os.Stdout).Encode(resp)
}
