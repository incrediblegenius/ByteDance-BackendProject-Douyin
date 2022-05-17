package client

import (
	"Douyin/proto/videoproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	VideoClient videoproto.VideosClient
)

func init() {
	conn, err := grpc.Dial("localhost:8889", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	VideoClient = videoproto.NewVideosClient(conn)
}
