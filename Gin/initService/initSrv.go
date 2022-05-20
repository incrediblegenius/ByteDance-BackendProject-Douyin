package initservice

import (
	"Douyin/global"
	"Douyin/proto"
	"Douyin/resolver"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	UserSrvConn()
}

func UserSrvConn() {

	conn, err := grpc.Dial(
		resolver.Target(fmt.Sprintf("http://nacos:nacos@%s:%d/nacos", global.NacosConfig.Host, global.NacosConfig.Port),
			global.ServerConfig.SrvServerInfo.UserSrv,
			resolver.OptionNameSpaceID(global.NacosConfig.Namespace),
			resolver.OptionGroupName(global.NacosConfig.Group),
			resolver.OptionNameSpaceID(global.NacosConfig.Namespace),
			resolver.OptionModeSubscribe()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		panic(err)
	}
	global.UserSrv = proto.NewServerClient(conn)

}
