package global

import (
	"UserServer/cfg"
	"UserServer/middleware"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/namsral/flag"
	"github.com/tencentyun/cos-go-sdk-v5"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	MysqlInfo   cfg.MysqlConfig
	DB          *gorm.DB
	Jwt         *middleware.JWT
	ServicePort int
	OssUrl      string
	OssClient   *cos.Client
	TempDir     string
	secretID    string
	secretKey   string
)

func init() {
	InitParse()
	os.MkdirAll(TempDir, os.ModePerm)
	InitOssClient()
	InitDB()
	Jwt = middleware.NewJWT()
}

func InitOssClient() {
	u, _ := url.Parse(OssUrl)
	b := &cos.BaseURL{BucketURL: u}
	OssClient = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: secretID,
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: secretKey,
		},
	})
}

func InitParse() {
	path, _ := os.Getwd()
	flag.StringVar(&MysqlInfo.User, "mysql_user", "root", "mysql user")
	flag.StringVar(&MysqlInfo.Password, "mysql_password", "root", "mysql password")
	flag.StringVar(&MysqlInfo.Host, "mysql_host", "127.0.0.1", "mysql host")
	flag.IntVar(&MysqlInfo.Port, "mysql_port", 3306, "mysql port")
	flag.StringVar(&MysqlInfo.Name, "mysql_name", "douyin_user", "mysql name")
	flag.IntVar(&ServicePort, "service_port", 8080, "service port")
	flag.StringVar(&OssUrl, "oss_url", "https://doiuyin-1302721364.cos.ap-chengdu.myqcloud.com", "oss url")
	flag.StringVar(&TempDir, "temp_dir", path+"/temp", "temp dir")
	flag.StringVar(&secretID, "secret_id", "AKIDAIC1cy62k7HDwQfhU4PWO32xhGgtvlOp", "oss secret id")
	flag.StringVar(&secretKey, "secret_key", "GI7lCPRIxPfjcIl14vZ3MTN4ZqsgI0Xs", "oss secret key")
	flag.Parse()
}

func InitDB() {
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", MysqlInfo.User, MysqlInfo.Password, MysqlInfo.Host, MysqlInfo.Port, MysqlInfo.Name)
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
