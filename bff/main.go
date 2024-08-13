package main

import (
	"github.com/Linxhhh/LinInk/bff/app"
	"github.com/Linxhhh/LinInk/bff/ioc"
	"github.com/gin-gonic/gin"
)

func main() {
	webserver := InitWebServer()
	webserver.Run(":8080")
}

func InitWebServer() *gin.Engine {

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
