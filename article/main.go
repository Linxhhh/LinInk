package main

import (
	"context"
	"log"
	"net"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	pb "github.com/Linxhhh/LinInk/api/proto/article"
	"github.com/Linxhhh/LinInk/article/events"
	server "github.com/Linxhhh/LinInk/article/grpc"
	"github.com/Linxhhh/LinInk/article/ioc"
	"github.com/Linxhhh/LinInk/article/repository"
	"github.com/Linxhhh/LinInk/article/repository/cache"
	"github.com/Linxhhh/LinInk/article/repository/dao"
	"github.com/Linxhhh/LinInk/article/service"
	"google.golang.org/grpc"
)

func main() {
	closeFunc := ioc.InitOtel()

	// init article repository
	cmd := ioc.InitCache()
	cache := cache.NewArticleCache(cmd)

	// init article dao
	master, slaves := ioc.InitDB()
	dao := dao.NewArticleDAO(master, slaves)

	// init article repository
	repo := repository.NewArticleRepository(dao, cache)

	// init user and interaction service client
	etcdCli := ioc.InitEtcdClient()
	userCli := ioc.InitUserRpcClient(etcdCli) 
	feedCli := ioc.InitFeedRpcClient(etcdCli)
	interCli := ioc.InitInteractionRpcClient(etcdCli)

	// init sarama producer
	cli := ioc.InitSaramaClient()
	pdr := ioc.InitSyncProducer(cli)

	// init publish and read event producer
	publishPdr := events.NewArticlePublishEventProducer(pdr)
	readPdr := events.NewArticleReadEventProducer(pdr)
	syncPdr := events.NewArticleSyncEventProducer(pdr)

	// init read event consumer
	readCsmr := events.NewArticleReadEventConsumer(cli, interCli)
	go func ()  {
		readCsmr.StartBatch("article_read")	
	}()

	// init article service
	svc := service.NewArticleService(repo, userCli, interCli, feedCli, publishPdr, readPdr, syncPdr)

	// init article service server
	svr := server.NewArticleServiceServer(svc)

	// create grpc server and register article service server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)
	pb.RegisterArticleServiceServer(grpcServer, svr)

	// Listen port
	listener, err := net.Listen("tcp", ":3336")
	if err != nil {
		log.Fatal(err)
	}

	// register to etcd
	if err = ioc.RegisterToEtcd(etcdCli); err != nil {
		log.Fatal(err)
	}

	// Start serve
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	closeFunc(ctx)
}