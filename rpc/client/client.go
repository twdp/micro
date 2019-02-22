package client

import (
	"github.com/hprose/hprose-golang/rpc"
)

var (
	client rpc.Client
	address []string
)

func init() {
	initAddress()
	client = rpc.NewClient(address...)
}

func UserService(remoteService interface{}, namespace ...string) {
	client.UseService(remoteService, namespace...)
}

func Close() {
	if client != nil {
		client.Close()
	}
}
