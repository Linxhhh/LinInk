package main

import (
	"log"
	"net"

	pb "github.com/Linxhhh/LinInk/api/proto/sms"
	server "github.com/Linxhhh/LinInk/sms/grpc"
	"github.com/Linxhhh/LinInk/sms/service"
	"google.golang.org/grpc"
)

func main() {
	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Create sms service server
	// svc := service.NewTentcentService()
	svc := service.NewLocalService()
	svr := server.NewSmsServiceServer(svc)

	// Register service
	pb.RegisterSmsServiceServer(grpcServer, svr)

	// Listen port
	listener, err := net.Listen("tcp", ":3335")
	if err != nil {
		log.Fatal(err)
	}

	// Start serve
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}