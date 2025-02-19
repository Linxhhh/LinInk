package main

import (
	"log"
	"net"

	pb "github.com/Linxhhh/LinInk/api/proto/comment"
	server "github.com/Linxhhh/LinInk/comment/grpc"
	"github.com/Linxhhh/LinInk/comment/ioc"
	"github.com/Linxhhh/LinInk/comment/repository"
	"github.com/Linxhhh/LinInk/comment/repository/dao"
	"github.com/Linxhhh/LinInk/comment/service"
	"google.golang.org/grpc"
)

func main() {

	m, s := ioc.InitDB()
	dao := dao.NewCommentDAO(m, s)
	repo := repository.NewCommentRepository(dao)

	etcdCli := ioc.InitEtcdClient()
	userCli := ioc.InitUserRpcClient(etcdCli)

	svc := service.NewCommentService(repo, userCli)
	svr := server.NewCommentServer(svc)

	// create grpc server and register article service server
	grpcServer := grpc.NewServer()
	pb.RegisterCommentServiceServer(grpcServer, svr)

	// Listen port
	listener, err := net.Listen("tcp", ":3332")
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
}
