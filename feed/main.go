package main

import (
	"log"
	"net"

	pb "github.com/Linxhhh/LinInk/api/proto/feed"
	"github.com/Linxhhh/LinInk/feed/events"
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
	etcdCli := ioc.InitEtcdClient()
	cli := ioc.InitFollowRpcClient(etcdCli)

	// init service
	svc := service.NewFeedEventService(repo, cli)

	// init sarama client
	saramaCli := ioc.InitSaramaClient()

	// init publish event consumer
	pubCsmr := events.NewArticlePublishEventConsumer(saramaCli, svc)
	go func ()  {
		pubCsmr.Start("article_publish")	
	}()

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

	// register to etcd
	if err = ioc.RegisterToEtcd(etcdCli); err != nil {
		log.Fatal(err)
	}

	// start serve
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}