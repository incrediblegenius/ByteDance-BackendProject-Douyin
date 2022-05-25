package main

import (
	"RelationSrv/global"
	"RelationSrv/handler"
	"RelationSrv/proto"
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
	fmt.Println("listen on:", global.ServicePort)
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	err = server.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}

}
