package global

import (
	"UserServer/cfg"
	"UserServer/middleware"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	NacosFileName = "config.yaml"
)

var (
	ServerConfig cfg.ServerConfig
	NacosConfig  cfg.NacosConfig
	DB           *gorm.DB
	Jwt          *middleware.JWT
	NamingClient naming_client.INamingClient
	ConfigClient config_client.IConfigClient
)

func init() {
	InitNacos()
	InitDB()
	Jwt = middleware.NewJWT()
}

func InitDB() {
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", ServerConfig.MysqlInfo.User, ServerConfig.MysqlInfo.Password, ServerConfig.MysqlInfo.Host, ServerConfig.MysqlInfo.Port, ServerConfig.MysqlInfo.Name)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	var err error
	DB, err = gorm.Open(mysql.Open(addr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("数据库连接成功")
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
