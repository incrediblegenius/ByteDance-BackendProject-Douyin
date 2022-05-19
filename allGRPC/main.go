package main

import (
	"UserServer/global"
	"UserServer/handler"
	"UserServer/proto"
	"UserServer/utils"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip address")
	port := flag.Int("port", 0, "port")
	if *port == 0 {
		*port, _ = utils.GetFreePort()
	}

	server := grpc.NewServer()
	proto.RegisterServerServer(server, &handler.Server{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	// 服务注册 nacos
	success, err := global.NamingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          *IP,
		Port:        uint64(*port),
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
	fmt.Println("服务注册成功，服务名称：", global.ServerConfig.Name, "，端口：", *port)
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// 注销服务
	success, err = global.NamingClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          *IP,
		Port:        uint64(*port),
		ServiceName: global.ServerConfig.Name,
		Ephemeral:   true,
		GroupName:   global.NacosConfig.Group, // 默认值DEFAULT_GROUP
	})
	if !success {
		panic(err)
	}
	fmt.Println("微服务注销成功")
}
