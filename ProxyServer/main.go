package main

import (
	"Douyin/ProxyServer/router"
	_ "Douyin/ProxyServer/userClient"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	router.InitRouter(r)
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
