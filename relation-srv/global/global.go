package global

import (
	"RelationSrv/cfg"
	"RelationSrv/middleware"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/namsral/flag"

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
)

func init() {
	InitParse()
	InitDB()
	Jwt = middleware.NewJWT()
}

func InitParse() {
	flag.StringVar(&MysqlInfo.User, "mysql_user", "root", "mysql user")
	flag.StringVar(&MysqlInfo.Password, "mysql_password", "root", "mysql password")
	flag.StringVar(&MysqlInfo.Host, "mysql_host", "127.0.0.1", "mysql host")
	flag.IntVar(&MysqlInfo.Port, "mysql_port", 3306, "mysql port")
	flag.StringVar(&MysqlInfo.Name, "mysql_name", "douyin_user", "mysql name")
	flag.IntVar(&ServicePort, "service_port", 8080, "service port")
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
