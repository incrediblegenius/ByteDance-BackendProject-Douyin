package router

import (
	"Douyin/ProxyServer/handler"

	"github.com/gin-gonic/gin"
)

func InitRouter(Router *gin.Engine) {
	r := Router.Group("/douyin/user")
	{
		r.POST("/register/", handler.Register)
		r.POST("/login/", handler.Login)
		r.GET("/", handler.GetUserInfo)

	}
	r = Router.Group("/douyin")
	{
		r.GET("/feed", handler.GetFeed)
	}
	r = Router.Group("/douyin/publish")
	{
		r.POST("/action/", handler.Publish)
		r.GET("/list/", handler.PublishList)
	}
}
