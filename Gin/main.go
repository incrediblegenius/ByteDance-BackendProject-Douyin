package main

import (
	"Douyin/global"
	_ "Douyin/initService"
	"Douyin/router"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	success, err := global.NamingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "0.0.0.0",
		Port:        8081,
		ServiceName: global.ServerConfig.Name,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    global.ServerConfig.Tags,
		ClusterName: "cluster-a",              // 默认值DEFAULT
		GroupName:   global.NacosConfig.Group, // 默认值DEFAULT_GROUP
	})
	if !success {
		panic(err)
	}
	glog.Info("register instance success")
	r := gin.Default()
	go func() {
		os.MkdirAll(global.ServerConfig.StaticInfo.StaticDir+"/videos", 0777)
		os.MkdirAll(global.ServerConfig.StaticInfo.StaticDir+"/covers", 0777)
		os.MkdirAll(global.ServerConfig.StaticInfo.StaticDir+"/tmp", 0777)
		http.Handle("/", http.FileServer(http.Dir(global.ServerConfig.StaticInfo.StaticDir)))
		if e := http.ListenAndServe(":8081", nil); e != nil {
			panic(e)
		}
	}()
	router.InitRouter(r)
	err = r.Run(":8080")
	if err != nil {
		panic(err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	success, err = global.NamingClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          "0.0.0.0",
		Port:        8081,
		ServiceName: global.ServerConfig.Name,
		Ephemeral:   true,
		GroupName:   global.NacosConfig.Group, // 默认值DEFAULT_GROUP
	})
	if !success {
		panic(err)
	}
	fmt.Println("微服务注销成功")
}
