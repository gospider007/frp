package main

import (
	"testing"

	"gitee.com/baixudong/frp"
)

func TestFprc(t *testing.T) {
	frp.NewClient(nil, frp.ClientOption{
		ServerHost: "111.111.111.111",
		ServerPort: 111,
		Token:      "111",

		Host: "127.0.0.1",
		Port: 111,

		Group:      "111",
		RemotePort: 111,
	})
}
