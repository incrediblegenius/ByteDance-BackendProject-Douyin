package grpcserver

import (
	"Douyin/proto/userproto"
	_ "Douyin/user_srv/global"
	"Douyin/user_srv/handler"
	"net"

	"google.golang.org/grpc"
)

func Run() error {
	server := grpc.NewServer()
	userproto.RegisterUserRegisterServer(server, &handler.UserRegisterServer{})
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