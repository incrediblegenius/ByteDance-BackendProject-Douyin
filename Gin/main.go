package main

import (
	"Douyin/global"
	_ "Douyin/global"
	"Douyin/router"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)
	fmt.Println("server start on port:", global.ServicePort)
	err := r.Run(fmt.Sprintf(":%d", global.ServicePort))
	if err != nil {
		panic(err)
	}
}
