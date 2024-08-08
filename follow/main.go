package main

import (
	"net"

	pb "github.com/Linxhhh/LinInk/api/proto/follow"
	server "github.com/Linxhhh/LinInk/follow/grpc"
	"github.com/Linxhhh/LinInk/follow/ioc"
	"github.com/Linxhhh/LinInk/follow/repository"
	"github.com/Linxhhh/LinInk/follow/repository/cache"
	"github.com/Linxhhh/LinInk/follow/repository/dao"
	"github.com/Linxhhh/LinInk/follow/service"
	"google.golang.org/grpc"
)

func main() {

	// init dao
	master, slaves := ioc.InitDB()
	dao := dao.NewFollowDAO(master, slaves)

	// init cache
	cmd := ioc.InitCache()
	cache := cache.NewFollowCache(cmd)

	// init service
	repo := repository.NewFollowRepository(dao, cache)
	svc := service.NewFollowService(repo)

	// create service server
	svr := server.NewFollowServiceServer(svc)

	// create grpc server and register server
	grpcServer := grpc.NewServer()
	pb.RegisterFollowServiceServer(grpcServer, svr)

	// listen port
	listener, err := net.Listen("tcp", ":3338")
	if err != nil {
		panic(err)
	}

	// start serve
	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}