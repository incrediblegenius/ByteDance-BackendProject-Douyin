package client

import (
	"Douyin/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	UserSrvClient proto.UserRegisterClient
)

func init() {
	conn, err := grpc.Dial("localhost:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	UserSrvClient = proto.NewUserRegisterClient(conn)
}
