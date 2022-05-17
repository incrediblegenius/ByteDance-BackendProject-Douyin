package client

import (
	"Douyin/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	VideoClient proto.VideosClient
)

func init() {
	conn, err := grpc.Dial("localhost:8889", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	VideoClient = proto.NewVideosClient(conn)
}
