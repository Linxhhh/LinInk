package main

import (
	"log"
	"net"

	pb "github.com/Linxhhh/LinInk/api/proto/search"
	"github.com/Linxhhh/LinInk/search/events"
	server "github.com/Linxhhh/LinInk/search/grpc"
	"github.com/Linxhhh/LinInk/search/ioc"
	"github.com/Linxhhh/LinInk/search/repository"
	"github.com/Linxhhh/LinInk/search/repository/dao"
	"github.com/Linxhhh/LinInk/search/service"
	"google.golang.org/grpc"
)

func main() {

	// ioc
	esClient := ioc.InitESClient()
	etcdClient := ioc.InitEtcdClient()

	// dao
	artDAO := dao.NewArticleElasticDAO(esClient)
	userDAO := dao.NewUserElasticDAO(esClient)

	// repo
	artRepo := repository.NewArticleRepository(artDAO)
	userRepo := repository.NewUserRepository(userDAO)

	// svc
	syncSvc := service.NewSyncService(artRepo, userRepo)
	searchSvc := service.NewSearchService(artRepo, userRepo)

	// csmr
	cli := ioc.InitSaramaClient()
	artCsmr := events.NewArticleSyncEventConsumer(cli, syncSvc)
	userCsmr := events.NewUserSyncEventConsumer(cli, syncSvc)
	go func ()  {
		artCsmr.Start("article_sync")	
	}()
	go func ()  {
		userCsmr.Start("user_sync")	
	}()

	// svr
	svr := server.NewSearchServiceServer(syncSvc, searchSvc)

	// Register service
	grpcServer := grpc.NewServer()
	pb.RegisterSearchServiceServer(grpcServer, svr)

	// Listen port
	listener, err := net.Listen("tcp", ":3331")
	if err != nil {
		log.Fatal(err)
	}

	// register to etcd
	if err = ioc.RegisterToEtcd(etcdClient); err != nil {
		log.Fatal(err)
	}

	// Start serve
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}