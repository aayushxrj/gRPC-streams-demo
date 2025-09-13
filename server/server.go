package main

import (
	"log"
	"net"
	"time"

	"github.com/aayushxrj/gRPC-streaming-demo/proto/gen"
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
		log.Println("Sent number:",a)
		if err != nil {
			return err
		}
		a,b = b, a+b
		time.Sleep(time.Second)
	}
	return nil
}

func main (){

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