package server

import (
	"UserServer/proto"
	_ "UserServer/user_srv/global"
	"UserServer/user_srv/handler"
	"net"

	"google.golang.org/grpc"
)

func Run() error {
	server := grpc.NewServer()
	proto.RegisterUserRegisterServer(server, &handler.UserRegisterServer{})
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
