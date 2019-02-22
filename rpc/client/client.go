package client

import (
	"github.com/hprose/hprose-golang/rpc"
	"tianwei.pro/micro/conf"
)

var (
	client rpc.Client
)

func init() {
	// 多客户端时修改此处，生成多个client

	client = rpc.NewClient(conf.Conf.Strings("client.address")...)
}

func UserService(remoteService interface{}, namespace ...string) {
	client.UseService(remoteService, namespace...)
}

func Close() {
	if client != nil {
		client.Close()
	}
}
