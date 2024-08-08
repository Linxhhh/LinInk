package main

import (
	"log"
	"net"

	pb "github.com/Linxhhh/LinInk/api/proto/feed"
	server "github.com/Linxhhh/LinInk/feed/grpc"
	"github.com/Linxhhh/LinInk/feed/ioc"
	"github.com/Linxhhh/LinInk/feed/repository"
	"github.com/Linxhhh/LinInk/feed/repository/cache"
	"github.com/Linxhhh/LinInk/feed/repository/dao"
	"github.com/Linxhhh/LinInk/feed/service"
	"google.golang.org/grpc"
)

func main() {

	// init dao
	master, slaves := ioc.InitDB()
	pulldao := dao.NewFeedPullEventDAO(master, slaves)
	pushdao := dao.NewFeedPushEventDAO(master, slaves)

	// init cache
	cmd := ioc.InitCache()
	cache := cache.NewFeedEventCache(cmd)
	
	// init repository
	repo := repository.NewFeedEventRepo(pulldao, pushdao, cache)
	
	// init client
	cli := ioc.InitFollowRpcClient()

	// init service
	svc := service.NewFeedEventService(repo, cli)

	// init service server
	svr := server.NewFeedServiceServer(svc)

	// init grpc server and register service server
	grpcServer := grpc.NewServer()
	pb.RegisterFeedServiceServer(grpcServer, svr)

	// listen port
	listener, err := net.Listen("tcp", ":3339")
	if err != nil {
		log.Fatal(err)
	}

	// start serve
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}