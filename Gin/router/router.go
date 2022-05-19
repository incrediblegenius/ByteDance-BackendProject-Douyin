package router

import (
	"Douyin/handler"

	"github.com/gin-gonic/gin"
)

func InitRouter(Router *gin.Engine) {
	Router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
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
	r = Router.Group("/douyin/favorite")
	{
		r.POST("/action/", handler.FavoriteAction)
		r.GET("/list/", handler.FavoriteList)
	}
	r = Router.Group("/douyin/comment")
	{
		r.POST("/action/", handler.CommentAction)
		r.GET("/list/", handler.CommentList)
	}
	r = Router.Group("/douyin/relation")
	{
		r.POST("/action/", handler.RelationAction)
		r.GET("/follow/list/", handler.FollowList)
		r.GET("/follower/list/", handler.FollowerList)
	}
}
