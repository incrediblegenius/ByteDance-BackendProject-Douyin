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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	NacosFileName = "config.yaml"
)

var (
	ServerConfig cfg.ServerConfig
	NacosConfig  cfg.NacosConfig
	SrvClient    proto.ServerClient
	NamingClient naming_client.INamingClient
	ConfigClient config_client.IConfigClient
)

func init() {
	InitNacos()
	// InitSrvConn()

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

func SrvConn() proto.ServerClient {
	instance, err := NamingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: ServerConfig.SrvServerInfo.Name,
		GroupName:   NacosConfig.Group,
	})
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", instance.Ip, instance.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	SrvClient = proto.NewServerClient(conn)

	return SrvClient
}
