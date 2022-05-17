package client

import (
	"Douyin/proto/userproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	UserSrvClient userproto.UserRegisterClient
)

func init() {
	conn, err := grpc.Dial("localhost:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	UserSrvClient = userproto.NewUserRegisterClient(conn)
}
