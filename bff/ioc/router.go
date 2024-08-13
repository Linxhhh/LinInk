package ioc

import (
	"github.com/Linxhhh/LinInk/bff/app"
	"github.com/gin-gonic/gin"
)

func InitRouter(halFunc []gin.HandlerFunc, userHdl *app.UserHandler, artHdl *app.ArticleHandler, followHdl *app.FollowHandler) *gin.Engine {
	router := gin.Default()
	router.Use(halFunc...)
	userHdl.RegistryRouter(router)
	artHdl.RegistryRouter(router)
	followHdl.RegistryRouter(router)
	return router
}