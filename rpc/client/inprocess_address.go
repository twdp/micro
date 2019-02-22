// +build !rpc

package client

import (
	"tianwei.pro/micro/conf"
	"tianwei.pro/micro/rpc/server"
)

func initAddress() {
	// 多客户端时修改此处，生成多个client
	address = conf.Conf.Strings("client.address")
	if len(address) == 0 {
		address = []string{
			server.Server.URI(),
		}
	}
}