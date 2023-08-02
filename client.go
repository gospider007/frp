package frp

import (
	"context"
	"errors"

	"gitee.com/baixudong/tools"
	frpc "github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/auth"
	"github.com/fatedier/frp/pkg/config"
	"github.com/fatedier/frp/pkg/util/log"
)

type Client struct {
	svr *frpc.Service
	ctx context.Context
}

func (obj *Client) Run() {
	obj.svr.Run(obj.ctx)
}
func (obj *Client) Close() {
	obj.svr.Close()
}

type ClientOption struct {
	ServerHost string //服务端host,默认0.0.0.0
	ServerPort int    //服务端port
	RemotePort int    //远程开放端口
	Host       string //本地服务host,默认0.0.0.0
	Port       int    //本地服务port
	Token      string //密钥，客户端与服务端连接验证
	Group      string // 负载均衡,分组
}

func NewClient(ctx context.Context, option ClientOption) (*Client, error) {
	if ctx == nil {
		ctx = context.TODO()
	}
	log.InitLog("console", "console", "error", 3, false)
	if option.Token == "" {
		return nil, errors.New("没有token,我认为你铁定连接不上服务")
	}
	if option.ServerHost == "" {
		option.ServerHost = "0.0.0.0"
	}
	if option.Host == "" {
		option.Host = "0.0.0.0"
	}
	if option.ServerPort == 0 {
		return nil, errors.New("没有设置监听端口,你确定能连接上服务")
	}
	if option.Port == 0 {
		return nil, errors.New("没有设置监听端口,你要转发到哪？")
	}
	if option.RemotePort == 0 {
		return nil, errors.New("没有设置开放端口,你要从哪接收外部流量？")
	}
	Name := tools.Uuid().String()
	svr, err := frpc.NewService(
		config.ClientCommonConf{
			ClientConfig: auth.ClientConfig{
				BaseConfig: auth.BaseConfig{
					AuthenticationMethod: "token",
				},
				TokenConfig: auth.TokenConfig{Token: option.Token},
			},
			Protocol:   "tcp",
			ServerAddr: option.ServerHost,
			ServerPort: option.ServerPort,

			DialServerTimeout:       10,
			DialServerKeepAlive:     7200,
			PoolCount:               1,
			TCPMux:                  true,
			TCPMuxKeepaliveInterval: 60,
			LoginFailExit:           true,
			QUICKeepalivePeriod:     10,
			QUICMaxIdleTimeout:      30,
			QUICMaxIncomingStreams:  100000,
			HeartbeatInterval:       30,
			HeartbeatTimeout:        90,
			UDPPacketSize:           1500,
		},
		map[string]config.ProxyConf{
			Name: &config.TCPProxyConf{
				RemotePort: option.RemotePort,
				BaseProxyConf: config.BaseProxyConf{
					Group:          option.Group,
					ProxyName:      Name,
					UseCompression: true,
					LocalSvrConf: config.LocalSvrConf{
						LocalIP:   option.Host,
						LocalPort: option.Port,
					},
				},
			},
		}, nil, "",
	)
	return &Client{svr: svr, ctx: ctx}, err
}
