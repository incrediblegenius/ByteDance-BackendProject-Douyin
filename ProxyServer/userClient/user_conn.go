package userClient

import (
	"Douyin/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	UserSrvClient proto.UserRegisterClient
)

func init() {
	conn, _ := grpc.Dial("localhost:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	UserSrvClient = proto.NewUserRegisterClient(conn)
}
