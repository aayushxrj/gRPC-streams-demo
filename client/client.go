package main

import (
	"context"
	"io"
	"log"
	"time"

	mainpb "github.com/aayushxrj/gRPC-streaming-demo/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	port := ":50051"
	cert := "cert.pem"

	// with TLS
	creds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		log.Fatal("Failed to load credentials:", err)
	}

	conn, err := grpc.NewClient("localhost"+port, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Println("Unable to connet", err)
	}
	defer conn.Close()

	client := mainpb.NewCalculatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &mainpb.FibonacciRequest{
		N: 10,
	}

	stream, err := client.GenerateFibonacci(ctx, req)
	if err != nil {
		log.Fatalln("Error calling GenerateFibonacci RPC", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Println("End of strearm data")
		}
		if err != nil {
			log.Fatalln("Error receiving data from RPC", err)
		}

		log.Println("Received Number:", res.Number)
	}
}
