package server

import (
	"github.com/hprose/hprose-golang/rpc"
	"tianwei.pro/micro/conf"
)

var (
	Server *rpc.TCPServer
)

func init() {
	Server = rpc.NewTCPServer(conf.Conf.String("server.address"))
	Handle()
}

func AddInstanceMethods(obj interface{}, option ...rpc.Options) {
	Server.AddInstanceMethods(obj, option...)
}

func Remove(name string) {
	Server.Remove(name)
}

func AddFunction(name string, function interface{}, option ...rpc.Options) {
	Server.AddFunction(name, function, option...)
}

func Close() {
	if Server != nil {
		Server.Close()
	}
}
