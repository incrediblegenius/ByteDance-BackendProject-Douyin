package global

import (
	"Douyin/proto"
	"context"

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
	localAddr       string
)

func InitParse() {
	flag.IntVar(&ServicePort, "service_port", 8080, "service port")
	flag.StringVar(&userService, "user_service_resolver", "127.0.0.1:8081", "user service resolver")
	flag.StringVar(&relationService, "relation_service_resolver", "127.0.0.1:8082", "relation service resolver")
	flag.StringVar(&feedService, "feed_service_resolver", "127.0.0.1:8083", "feed service resolver")
	flag.StringVar(&commentService, "comment_service_resolver", "127.0.0.1:8084", "comment service resolver")
	flag.StringVar(&favoriteService, "favorite_service_resolver", "127.0.0.1:8085", "favorite service resolver")
	flag.StringVar(&publishService, "publish_service_resolver", "127.0.0.1:8086", "publish service resolver")
	flag.StringVar(&localAddr, "local_addr", "", "local addr")
	flag.Parse()
}

func InitSrvClient() {
	if localAddr != "" {
		userService = localAddr + ":8081"
		relationService = localAddr + ":8082"
		feedService = localAddr + ":8083"
		commentService = localAddr + ":8084"
		favoriteService = localAddr + ":8085"
		publishService = localAddr + ":8086"
	}
	var conn *grpc.ClientConn
	var err error
	if conn, err = grpc.DialContext(
		context.Background(),
		userService,
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
		relationService,
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
		feedService,
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
		commentService,
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
		favoriteService,
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
		publishService,
		// "127.0.0.1:8080",
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
	InitParse()
	InitSrvClient()
}
