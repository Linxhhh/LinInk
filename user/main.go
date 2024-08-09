package main

import (
	"log"
	"net"

	pb "github.com/Linxhhh/LinInk/api/proto/user"
	server "github.com/Linxhhh/LinInk/user/grpc"
	"github.com/Linxhhh/LinInk/user/ioc"
	"github.com/Linxhhh/LinInk/user/repository"
	"github.com/Linxhhh/LinInk/user/repository/cache"
	"github.com/Linxhhh/LinInk/user/repository/dao"
	"github.com/Linxhhh/LinInk/user/service"
	"google.golang.org/grpc"
)

func main() {
	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Init code cache
	cmd := ioc.InitCache()
	cache := cache.NewUserCache(cmd)

	// Init code DAO
	master, slave := ioc.InitDB()
	dao := dao.NewUserDAO(master, slave)
	
	// Create code service server
	repo := repository.NewUserRepository(dao, cache)
	svc := service.NewUserService(repo)
	svr := server.NewUserServiceServer(svc)

	// Register service
	pb.RegisterUserServiceServer(grpcServer, svr)

	// Listen port
	listener, err := net.Listen("tcp", ":3333")
	if err != nil {
		log.Fatal(err)
	}

	// register to etcd
	etcdCli := ioc.InitEtcdClient()
	if err = ioc.RegisterToEtcd(etcdCli); err != nil {
		log.Fatal(err)
	}

	// Start serve
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}