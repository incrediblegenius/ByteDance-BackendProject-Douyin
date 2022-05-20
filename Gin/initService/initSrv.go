package initservice

import (
	"Douyin/global"
	"Douyin/proto"
	"Douyin/resolver"
	"fmt"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	value := reflect.ValueOf(global.ServerConfig.SrvServerInfo)
	fmt.Println(global.ServerConfig.SrvServerInfo)
	for num := 0; num < value.NumField(); num++ {
		n := value.Field(num).String()
		if n != "" {
			SrvConn(n)
		}
	}

}

func SrvConn(srv string) {
	conn, err := grpc.Dial(
		resolver.Target(fmt.Sprintf("http://nacos:nacos@%s:%d/nacos", global.NacosConfig.Host, global.NacosConfig.Port),
			srv,
			resolver.OptionNameSpaceID(global.NacosConfig.Namespace),
			resolver.OptionGroupName(global.NacosConfig.Group),
			resolver.OptionModeSubscribe()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		panic(err)
	}
	global.ConnMap[srv] = proto.NewServerClient(conn)
}
