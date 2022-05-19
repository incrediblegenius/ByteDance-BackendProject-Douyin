package cfg

type StaticServerConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	StaticDir string `json:"static_dir"`
}

type SrvServerConfig struct {
	Name string `json:"name"`
}

type ServerConfig struct {
	Name          string             `json:"name"`
	Tags          map[string]string  `json:"tag"`
	StaticInfo    StaticServerConfig `json:"static_server"`
	SrvServerInfo SrvServerConfig    `json:"srv_server"`
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
