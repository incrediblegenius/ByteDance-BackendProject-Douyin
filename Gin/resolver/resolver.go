package resolver

// Nacos 的grpc解析器（参考github.com/magicdvd/nacos-grpc）

import (
	"Douyin/global"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc/resolver"
)

func init() {
	resolver.Register(NewBuilder())
}

const (
	modeHeartBeat             string = "hb"
	modeSubscribe             string = "sb"
	DefaultClusterName               = "cluster-a"
	DefaultGroupName                 = "DEFAULT_GROUP"
	DefaultNameSpaceID               = "public"
	DefaultTimeout                   = 15 * time.Second
	DefaultListenInterval            = 30 * time.Second
	DefaultMaxCacheTime              = 45 * time.Second
	DefaultSubscrubeCacheTime        = 10 * time.Second
)

var (
	ErrUnsupportSchema = errors.New("unsupport schema (nacos/nacoss)")
	ErrMissServiceName = errors.New("target miss service name")
	ErrMissGroupName   = errors.New("target miss group name")
	ErrMissNameSpaceID = errors.New("target miss namespace name")
	ErrMode            = errors.New("target mode err (hb/sb)")
	ErrMissInterval    = errors.New("target mode heartbeat miss interval")
	ErrInterval        = errors.New("target mode heartbeat interval error")
	ErrNoInstances     = errors.New("no valid instance")
)

type Option interface {
	apply(opts *options)
}

type options struct {
	groupName   string
	clusters    string
	nameSpaceID string
	mode        string
	hbInterval  time.Duration
}

type op struct {
	f func(opts *options)
}

func (c *op) apply(opts *options) {
	c.f(opts)
}

func OptionGroupName(g string) Option {
	return &op{
		f: func(opts *options) {
			opts.groupName = g
		},
	}
}

func OptionNameSpaceID(g string) Option {
	return &op{
		f: func(opts *options) {
			opts.nameSpaceID = g
		},
	}
}

func OptionClusters(c []string) Option {
	return &op{
		f: func(opts *options) {
			if len(c) > 0 {
				opts.clusters = strings.Join(c, ",")
			}
		},
	}
}

func OptionModeHeartBeat(d time.Duration) Option {
	return &op{
		f: func(opts *options) {
			opts.mode = modeHeartBeat
			opts.hbInterval = d
		},
	}
}

func OptionModeSubscribe() Option {
	return &op{
		f: func(opts *options) {
			opts.mode = modeSubscribe
		},
	}
}

func Target(nacosAddr string, serviceName string, ops ...Option) string {
	opts := &options{
		groupName:   DefaultGroupName,
		clusters:    DefaultClusterName,
		nameSpaceID: DefaultNameSpaceID,
		mode:        modeSubscribe,
		hbInterval:  20 * time.Second,
	}
	for _, v := range ops {
		v.apply(opts)
	}
	tmp := nacosAddr
	if strings.HasPrefix(nacosAddr, "https://") {
		tmp = "nacoss://" + nacosAddr[8:]
	} else if strings.HasPrefix(nacosAddr, "http://") {
		tmp = "nacos://" + nacosAddr[7:]
	}
	return fmt.Sprintf("%s?s=%s&n=%s&cs=%s&g=%s&m=%s&d=%d", tmp, serviceName, opts.nameSpaceID, opts.clusters, opts.groupName, opts.mode, opts.hbInterval/time.Millisecond)
}

func NewBuilder() resolver.Builder {
	return &nacosResolverBuilder{}
}

type nacosResolverBuilder struct{}

// // Build
func (*nacosResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r, err := newNacosResolver(target, cc)
	if err != nil {
		return nil, err
	}
	go r.start()
	return r, nil
}

func (*nacosResolverBuilder) Scheme() string {
	return "nacos"
}

type nacosResolver struct {
	namingClient naming_client.INamingClient
	cc           resolver.ClientConn
	mode         string
	close        chan bool
	interval     time.Duration
	serviceName  string
	groupName    string
	nameSpaceId  string
	clusters     string
}

func newNacosResolver(target resolver.Target, cc resolver.ClientConn) (*nacosResolver, error) {
	if target.URL.Scheme != "nacos" && target.URL.Scheme != "nacoss" {
		return nil, ErrUnsupportSchema
	}
	u, err := url.Parse("http://test.com/?" + target.URL.RawQuery)
	if err != nil {
		return nil, err
	}

	values := u.Query()

	serviceName := values.Get("s")
	if serviceName == "" {
		return nil, ErrMissServiceName
	}
	nameSpaceId := values.Get("n")
	if nameSpaceId == "" {
		return nil, ErrMissNameSpaceID
	}
	groupName := values.Get("g")
	if groupName == "" {
		return nil, ErrMissGroupName
	}
	clusterName := values.Get("cs")
	if clusterName == "" {
		return nil, ErrMissGroupName
	}
	mode := values.Get("m")
	if mode != modeHeartBeat && mode != modeSubscribe {
		return nil, ErrMode
	}
	var interval time.Duration
	if mode == modeHeartBeat {
		s := values.Get("d")
		if s == "" && mode == modeHeartBeat {
			return nil, ErrMissInterval
		}
		tmp, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		if tmp <= 0 {
			return nil, ErrInterval
		}
		interval = time.Duration(tmp) * time.Millisecond
	}
	c := &nacosResolver{
		namingClient: global.NamingClient,
		cc:           cc,
		mode:         mode,
		close:        make(chan bool),
		interval:     interval,
		serviceName:  serviceName,
		groupName:    groupName,
		nameSpaceId:  nameSpaceId,
		clusters:     clusterName,
	}
	return c, nil
}

func (c *nacosResolver) start() {
	addrs, err := c.getInstances()
	if err != nil {
		c.cc.ReportError(err)
	} else {
		c.cc.UpdateState(resolver.State{Addresses: addrs})
	}
	if c.mode == modeHeartBeat {
		tick := time.NewTicker(c.interval)
		for {
			select {
			case <-tick.C:
				addrs, err := c.getInstances()
				if err != nil {
					c.cc.ReportError(err)
				} else {
					c.cc.UpdateState(resolver.State{Addresses: addrs})
				}
			case <-c.close:
				return
			}
		}
	} else {
		err := c.namingClient.Subscribe(&vo.SubscribeParam{
			ServiceName: c.serviceName,
			GroupName:   c.groupName,           // 默认值DEFAULT_GROUP
			Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
			SubscribeCallback: func(services []model.SubscribeService, e error) {
				addrs, err := c.getInstances()
				if err != nil {
					c.cc.ReportError(err)
				} else {
					c.cc.UpdateState(resolver.State{Addresses: addrs})
				}
			},
		})
		if err != nil {
			c.cc.ReportError(err)
		}
	}
}

func (c *nacosResolver) getInstances() ([]resolver.Address, error) {
	addrs, err := c.namingClient.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: c.serviceName,
		GroupName:   c.groupName,          // 默认值DEFAULT_GROUP
		Clusters:    []string{c.clusters}, // 默认值DEFAULT
	})
	if err != nil || len(addrs) == 0 {
		return nil, ErrNoInstances
	}
	l := len(addrs)
	ret := make([]resolver.Address, l)
	for i := 0; i < l; i++ {
		if !addrs[i].Healthy {
			continue
		}
		addr := resolver.Address{
			Addr:       fmt.Sprintf("%s:%d", addrs[i].Ip, addrs[i].Port),
			ServerName: c.serviceName,
		}
		ret = append(ret, addr)
	}
	return ret, nil
}

func (c *nacosResolver) ResolveNow(o resolver.ResolveNowOptions) {
	//directly get service
}

func (c *nacosResolver) Close() {
	if c.mode == modeHeartBeat {
		close(c.close)
	} else {
		c.namingClient.Unsubscribe(&vo.SubscribeParam{
			ServiceName: c.serviceName,
			GroupName:   c.groupName,          // 默认值DEFAULT_GROUP
			Clusters:    []string{c.clusters}, // 默认值DEFAULT
			SubscribeCallback: func(services []model.SubscribeService, err error) {
				log.Printf("服务发现监听注销")
			},
		})
	}
}
