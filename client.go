package frp

import (
	"context"
	"errors"

	frpc "github.com/fatedier/frp/client"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/util/log"
	"github.com/google/uuid"
)

type Client struct {
	svr *frpc.Service
	ctx context.Context
}

func (obj *Client) Run() error {
	return obj.svr.Run(obj.ctx)
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
	log.InitLogger("console", "error", 3, false)
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
	if ctx == nil {
		ctx = context.TODO()
	}
	svr, err := frpc.NewService(
		frpc.ServiceOptions{
			Common: &v1.ClientCommonConfig{
				Auth: v1.AuthClientConfig{
					Method: v1.AuthMethodToken,
					Token:  option.Token,
				},
				ServerAddr: option.ServerHost,
				ServerPort: option.ServerPort,
			},
			ProxyCfgs: []v1.ProxyConfigurer{
				&v1.TCPProxyConfig{
					RemotePort: option.RemotePort,
					ProxyBaseConfig: v1.ProxyBaseConfig{
						ProxyBackend: v1.ProxyBackend{
							LocalIP:   option.Host,
							LocalPort: option.Port,
						},
						Transport: v1.ProxyTransport{
							UseCompression: true,
						},
						LoadBalancer: v1.LoadBalancerConfig{
							Group: option.Group,
						},
						Name: uuid.New().String(),
					},
				},
			},
		},
	)
	return &Client{svr: svr, ctx: ctx}, err
}
