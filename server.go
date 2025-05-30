package frp

import (
	"context"
	"errors"

	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/util/log"
	frps "github.com/fatedier/frp/server"
)

type Server struct {
	svr *frps.Service
	ctx context.Context
}

func (obj *Server) Run() {
	obj.svr.Run(obj.ctx)
}

type ServerOption struct {
	Host  string //服务端host,默认0.0.0.0
	Port  int    //服务端port
	Token string //密钥，客户端与服务端连接验证
}

func NewServer(ctx context.Context, option ServerOption) (*Server, error) {
	if ctx == nil {
		ctx = context.TODO()
	}
	log.InitLogger("console", "error", 3, false)
	if option.Token == "" {
		return nil, errors.New("没有token,你想被攻击吗？")
	}
	if option.Host == "" {
		option.Host = "0.0.0.0"
	}
	if option.Port == 0 {
		return nil, errors.New("服务端没有设置监听端口,你确定要这样？")
	}
	svr, err := frps.NewService(
		&v1.ServerConfig{
			Auth: v1.AuthServerConfig{
				Method: v1.AuthMethodToken,
				Token:  option.Token,
			},
			BindAddr: option.Host,
			BindPort: option.Port,
		},
	)
	return &Server{svr: svr, ctx: ctx}, err
}
