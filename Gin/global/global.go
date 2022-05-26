package global

import (
	"Douyin/proto"
	"context"
	"fmt"

	"github.com/namsral/flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ConnMap         = make(map[string]proto.ServerClient)
	ServicePort     int
	userService     string
	relationService string
	feedService     string
	commentService  string
	favoriteService string
	publishService  string
)

func InitParse() {
	flag.IntVar(&ServicePort, "service_port", 8080, "service port")
	flag.StringVar(&userService, "user_service", "user-balance-svc", "headless-service of user")
	flag.StringVar(&relationService, "relation_service", "relation-balance-svc", "headless-service of relation")
	flag.StringVar(&feedService, "feed_service", "feed-balance-svc", "headless-service of feed")
	flag.StringVar(&commentService, "comment_service", "comment-balance-svc", "headless-service of comment")
	flag.StringVar(&favoriteService, "favorite_service", "favorite-balance-svc", "headless-service of favorite")
	flag.StringVar(&publishService, "publish_service", "publish-balance-svc", "headless-service of publish")
	flag.Parse()
}

func InitSrvClient() {
	var conn *grpc.ClientConn
	var err error
	if conn, err = grpc.DialContext(
		context.Background(),
		fmt.Sprintf("dns:///%s", userService),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		// grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*10)),
		grpc.WithBlock(),
	); err != nil {
		panic(err)
	}
	ConnMap["user_srv"] = proto.NewServerClient(conn)

	if conn, err = grpc.DialContext(
		context.Background(),
		fmt.Sprintf("dns:///%s", relationService),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		// grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*10)),
		grpc.WithBlock(),
	); err != nil {
		panic(err)
	}
	ConnMap["relation_srv"] = proto.NewServerClient(conn)

	if conn, err = grpc.DialContext(
		context.Background(),
		fmt.Sprintf("dns:///%s", feedService),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		// grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*10)),
		grpc.WithBlock(),
	); err != nil {
		panic(err)
	}
	ConnMap["feed_srv"] = proto.NewServerClient(conn)

	if conn, err = grpc.DialContext(
		context.Background(),
		fmt.Sprintf("dns:///%s", commentService),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		// grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*10)),
		grpc.WithBlock(),
	); err != nil {
		panic(err)
	}
	ConnMap["comment_srv"] = proto.NewServerClient(conn)

	if conn, err = grpc.DialContext(
		context.Background(),
		fmt.Sprintf("dns:///%s", favoriteService),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		// grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*10)),
		grpc.WithBlock(),
	); err != nil {
		panic(err)
	}
	ConnMap["favorite_srv"] = proto.NewServerClient(conn)

	if conn, err = grpc.DialContext(
		context.Background(),
		fmt.Sprintf("dns:///%s", publishService),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*10)),
		grpc.WithBlock(),
	); err != nil {
		panic(err)
	}
	ConnMap["publish_srv"] = proto.NewServerClient(conn)

}

func init() {

}
