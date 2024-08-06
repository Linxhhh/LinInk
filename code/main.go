package main

import (
	"log"
	"net"

	pb "github.com/Linxhhh/LinInk/api/proto/code"
	server "github.com/Linxhhh/LinInk/code/grpc"
	"github.com/Linxhhh/LinInk/code/ioc"
	"github.com/Linxhhh/LinInk/code/repository"
	"github.com/Linxhhh/LinInk/code/repository/cache"
	"github.com/Linxhhh/LinInk/code/service"
	"google.golang.org/grpc"
)

func main() {
	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Init code repository
	cmd := ioc.InitCache()
	cache := cache.NewCodeCache(cmd)
	repo := repository.NewCodeRepository(cache)
	
	// Create sms service client
	cli := ioc.InitSmsRpcClient()

	// Create code service server
	svc := service.NewCodeService(repo, cli)
	svr := server.NewCodeServiceServer(svc)

	// Register service
	pb.RegisterCodeServiceServer(grpcServer, svr)

	// Listen port
	listener, err := net.Listen("tcp", ":3334")
	if err != nil {
		log.Fatal(err)
	}

	// Start serve
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}