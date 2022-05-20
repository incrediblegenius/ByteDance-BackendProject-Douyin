package cfg

type StaticServerConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	StaticDir string `json:"static_dir"`
}

type ServerConfig struct {
	Name          string             `json:"name"`
	Tags          map[string]string  `json:"tag"`
	StaticInfo    StaticServerConfig `json:"static_server"`
	SrvServerInfo ServiceConfig      `json:"srv_server"`
}

type ServiceConfig struct {
	UserSrv     string `json:"user_srv"`
	FeedSrv     string `json:"feed_srv"`
	PublishSrv  string `json:"publish_srv"`
	FavoriteSrv string `json:"favorite_srv"`
	RelationSrv string `json:"relation_srv"`
	CommentSrv  string `json:"comment_srv"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
