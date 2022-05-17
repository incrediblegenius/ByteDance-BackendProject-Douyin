package main

import (
	_ "Douyin/ProxyServer/client"
	"Douyin/ProxyServer/router"

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