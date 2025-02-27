package main

import (
	"log"
	"net"

	pb "github.com/Linxhhh/LinInk/api/proto/interaction"
	"github.com/Linxhhh/LinInk/interaction/events"
	server "github.com/Linxhhh/LinInk/interaction/grpc"
	"github.com/Linxhhh/LinInk/interaction/ioc"
	"github.com/Linxhhh/LinInk/interaction/repository"
	"github.com/Linxhhh/LinInk/interaction/repository/cache"
	"github.com/Linxhhh/LinInk/interaction/repository/dao"
	"github.com/Linxhhh/LinInk/interaction/service"
	"google.golang.org/grpc"
)

func main() {

	// init interaction dao
	master, slaves := ioc.InitDB()
	dao := dao.NewInteractionDAO(master, slaves)

	// init interaction cache
	cmd := ioc.InitCache()
	cache := cache.NewInteractionCache(cmd)

	// init interaction repository
	repo := repository.NewInteractionRepository(dao, cache)

	// init interaction service
	svc := service.NewInteractionService(repo)

	// init sarama client
	saramaCli := ioc.InitSaramaClient()

	// init read event consumer
	readCsmr := events.NewArticleReadEventConsumer(saramaCli, svc)
	go func ()  {
		readCsmr.StartBatch("article_read")	
	}()

	// Create interaction service server
	svr := server.NewInteractionServiceServer(svc)

	// Create grpc server and register interaction service server
	grpcServer := grpc.NewServer()
	pb.RegisterInteractionServiceServer(grpcServer, svr)

	// Listen port
	listener, err := net.Listen("tcp", ":3337")
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
