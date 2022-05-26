package main

import (
	"UserServer/global"
	"UserServer/handler"
	"UserServer/proto"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {

	server := grpc.NewServer()
	proto.RegisterServerServer(server, &handler.Server{})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", global.ServicePort))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	fmt.Println("grpc server start at port:", global.ServicePort)
	err = server.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
