package server

import (
	"github.com/hprose/hprose-golang/rpc"
	"tianwei.pro/micro/conf"
)

var (
	server *rpc.TCPServer
)

func init() {
	server = rpc.NewTCPServer(conf.Conf.String("server.address"))
	server.Handle()
}

func AddInstanceMethods(obj interface{}, option ...rpc.Options) {
	server.AddInstanceMethods(obj, option...)
}

func Remove(name string) {
	server.Remove(name)
}

func AddFunction(name string, function interface{}, option ...rpc.Options) {
	server.AddFunction(name, function, option...)
}

func Close() {
	if server != nil {
		server.Close()
	}
}
