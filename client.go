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
	Host       string //本地服务host,默认0.0.0.0
	Port       int    //本地服务port
	RemotePort int    //远程开放端口
	Group      string // 负载均衡,分组
}

func NewClient(ctx context.Context, serverOption ServerOption, clientOption ...ClientOption) (*Client, error) {
	if len(clientOption) == 0 {
		return nil, errors.New("ErrNoClientOption")
	}
	log.InitLogger("console", "error", 3, false)
	if ctx == nil {
		ctx = context.TODO()
	}
	ProxyCfgs := []v1.ProxyConfigurer{}
	for _, option := range clientOption {
		ProxyCfgs = append(ProxyCfgs, &v1.TCPProxyConfig{
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
		)
	}
	svr, err := frpc.NewService(
		frpc.ServiceOptions{
			Common: &v1.ClientCommonConfig{
				Auth: v1.AuthClientConfig{
					Method: v1.AuthMethodToken,
					Token:  serverOption.Token,
				},
				ServerAddr: serverOption.Host,
				ServerPort: serverOption.Port,
			},
			ProxyCfgs: ProxyCfgs,
		},
	)
	return &Client{svr: svr, ctx: ctx}, err
}
