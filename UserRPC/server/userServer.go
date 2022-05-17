package server

import (
	_ "UserServer/global"
	"UserServer/handler"
	"UserServer/proto"
	"net"

	"google.golang.org/grpc"
)

func Run() error {
	server := grpc.NewServer()
	proto.RegisterServerServer(server, &handler.Server{})
	lis, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
	return nil
}
