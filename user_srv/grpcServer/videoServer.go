package grpcserver

import (
	"Douyin/proto"
	"Douyin/user_srv/handler"
	"net"

	"google.golang.org/grpc"
)

func VideoServerRun() error {
	server := grpc.NewServer()
	proto.RegisterVideosServer(server, &handler.VideosServer{})
	lis, err := net.Listen("tcp", "localhost:8889")
	if err != nil {
		panic(err)
	}
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
	return nil
}
