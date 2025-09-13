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

	// Server side streaming

	req := &mainpb.FibonacciRequest{
		N: 1,
	}

	stream, err := client.GenerateFibonacci(ctx, req)
	if err != nil {
		log.Fatalln("Error calling GenerateFibonacci RPC", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Println("End of strearm data")
			break
		}
		if err != nil {
			log.Println("Error receiving data from RPC", err)
		}

		log.Println("Received Number:", res.Number)
	}

	// Client side streaming

	stream1, err := client.SendNumbers(ctx)
	if err != nil {
		log.Fatalln("Error creating stream", err)
	}
	
	for num := range 10 {
		err := stream1.Send(&mainpb.NumberRequest{Number : int32(num)})
		if err != nil {
			log.Fatalln("Error sending Number:", err)
		}
	} 

	res, err := stream1.CloseAndRecv()
	if err != nil {
		log.Println("Error receiving response:", err)
	}

	log.Println("Sum is:",res.Sum)

	// Bi-directional streaming
	chatStream, err := client.Chat(ctx)
	if err != nil {
		log.Fatalln("Error creating chat stream:", err)
	}

	waitc := make(chan struct{})
	// Send messages in a goroutine
	go func() {
		messages := []string{"Hello", "How are you?", "Goodbye"}
		for _, message := range messages {
			log.Println("Sending message:", message)
			err := chatStream.Send(&mainpb.ChatMessage{Message: message})
			if err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Second)
		}
		chatStream.CloseSend()
	}()

	// Receive messages in goroutine
	go func() {
		for {
			res, err := chatStream.Recv()
			if err == io.EOF {
				log.Println("End of stream")
				break
			}
			if err != nil {
				log.Fatalln("Error receiving data from GenerateFibonacci func:", err)
			}
			log.Println("Received response: ", res.GetMessage())
		}
		close(waitc)
	}()
	<-waitc
}
