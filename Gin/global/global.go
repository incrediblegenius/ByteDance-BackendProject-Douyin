package global

import (
	"Douyin/cfg"
	"Douyin/proto"
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

const (
	NacosFileName = "config.yaml"
)

var (
	ServerConfig cfg.ServerConfig
	NacosConfig  cfg.NacosConfig

	NamingClient naming_client.INamingClient
	ConfigClient config_client.IConfigClient

	ConnMap = make(map[string]proto.ServerClient)
)

func init() {
	InitNacos()
}

func InitNacos() {
	v := viper.New()
	v.SetConfigFile(NacosFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&NacosConfig); err != nil {
		panic(err)
	}
	sc := []constant.ServerConfig{
		{
			IpAddr: NacosConfig.Host,
			Port:   NacosConfig.Port,
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "nacos/log",
		CacheDir:            "nacos/cache",
		LogLevel:            "warn",
	}
	var err error
	ConfigClient, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}
	NamingClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}
	content, err := ConfigClient.GetConfig(vo.ConfigParam{
		DataId: NacosConfig.DataId,
		Group:  NacosConfig.Group})
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(content), &ServerConfig)
	if err != nil {
		panic(err)
	}
	err = ConfigClient.ListenConfig(vo.ConfigParam{
		DataId: NacosConfig.DataId,
		Group:  NacosConfig.Group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("group:" + group + ", dataId:" + dataId + "配置发生改变！")
			err = json.Unmarshal([]byte(data), &ServerConfig)
		},
	})
	if err != nil {
		panic(err)
	}
	// fmt.Println(ServerConfig)
	fmt.Println("成功从nacos读取配置")
}

// func UserSrvConn() proto.ServerClient {
// 	// nacos 的负载均衡设置TODO
// 	// instance, err := NamingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
// 	// 	ServiceName: ServerConfig.SrvServerInfo.UserSrv,
// 	// 	GroupName:   NacosConfig.Group,
// 	// })
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	conn, err := grpc.Dial(
// 		nacosgrpc.Target(fmt.Sprintf("http://nacos:nacos@%s:%d/nacos", NacosConfig.Host, NacosConfig.Port),
// 			ServerConfig.SrvServerInfo.UserSrv,
// 			nacosgrpc.OptionNameSpaceID(NacosConfig.Namespace),
// 			nacosgrpc.OptionGroupName(NacosConfig.Group)),
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	UserSrv = proto.NewServerClient(conn)

// 	return UserSrv

// }
