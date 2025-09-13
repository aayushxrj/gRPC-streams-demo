package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	mainpb "github.com/aayushxrj/gRPC-streaming-demo/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	mainpb.UnimplementedCalculatorServer
}

func (s *server) GenerateFibonacci(req *mainpb.FibonacciRequest, stream mainpb.Calculator_GenerateFibonacciServer) error {
	n := req.GetN()
	a, b := 0, 1

	for i := 0; i < int(n); i++ {
		err := stream.Send(&mainpb.FibonacciResponse{
			Number: int32(a),
		})
		log.Println("Sent number:", a)
		if err != nil {
			return err
		}
		a, b = b, a+b
		time.Sleep(time.Second)
	}
	return nil
}

func (s *server) SendNumbers(stream mainpb.Calculator_SendNumbersServer) error {
	var sum int32

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&mainpb.NumberResponse{Sum: sum})
		}
		if err != nil {
			return err
		}
		log.Println(req.GetNumber())

		sum += req.GetNumber()
	}
}

func (s *server) Chat(stream mainpb.Calculator_ChatServer) error {

	// read from terminal
	reader := bufio.NewReader(os.Stdin)

	for {
		// receiving value/messages from stream
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Received Message:", req.GetMessage())

		// sending value/messages through the stream

		fmt.Print("Enter response:")
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		msg = strings.TrimSpace(msg)

		err = stream.Send(&mainpb.ChatMessage{
			Message: msg,
			// Message: req.GetMessage(),
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println("Returning control")
	return nil
}

func main() {

	port := ":50051"
	cert := "cert.pem"
	key := "key.pem"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		log.Fatal("Failed to load credentials:", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	mainpb.RegisterCalculatorServer(grpcServer, &server{})

	log.Printf("Server is running on the port%s", port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
