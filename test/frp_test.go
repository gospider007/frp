package main

import (
	"log"
	"testing"
	"time"

	"github.com/gospider007/frp"
)

func TestFprc(t *testing.T) {
	go func() {
		ff, err := frp.NewServer(nil, frp.ServerOption{
			Port:  111,
			Token: "111",
		})
		if err != nil {
			log.Panic(err)
		}
		ff.Run()
	}()
	time.Sleep(time.Second * 2)
	cli, err := frp.NewClient(nil, frp.ClientOption{
		// ServerHost: "111.111.111.111",
		ServerHost: "127.0.0.1",
		ServerPort: 111,
		Token:      "111",

		Host: "192.168.110.26",
		Port: 27017,

		Group:      "111",
		RemotePort: 27017,
	})
	if err != nil {
		log.Panic(err)
	}
	log.Print(cli.Run())
}
