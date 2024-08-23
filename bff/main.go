package main

import (
	"context"
	"net/http"
	"time"

	"github.com/Linxhhh/LinInk/bff/app"
	"github.com/Linxhhh/LinInk/bff/ioc"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	clostFunc := ioc.InitOtel()

	svr := initWebServer()
	initPrometheus()
	svr.Run(":8081")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	clostFunc(ctx)
}

func initWebServer() *gin.Engine {

	// client
	cli := ioc.InitEtcdClient()
	userCli := ioc.InitUserRpcClient(cli)
	codeCli := ioc.InitCodeRpcClient(cli)
	followCli := ioc.InitFollowRpcClient(cli)
	articleCli := ioc.InitArticleRpcClient(cli)
	interactionCli := ioc.InitInteractionRpcClient(cli)

	// handler
	userHandler := app.NewUserHandler(userCli, codeCli)
	articleHandler := app.NewArticleHandler(articleCli, interactionCli)
	followHandler := app.NewFollowHandler(followCli)

	// middleware
	hdlFuncs := ioc.InitMiddleware()

	// router
	return ioc.InitRouter(hdlFuncs, userHandler, articleHandler, followHandler)
}

func initPrometheus() {
	go func ()  {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8082", nil)
	}()
}