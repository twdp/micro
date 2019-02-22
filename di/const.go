package di

import "fmt"

const (
	RpcProvider = "rpc_s_"
	RpcConsumer = "rpc_c_"
)

func NewRpcProviderName(name string) string {
	return fmt.Sprintf("%s%s", RpcProvider, name)
}

func NewRpcConsumerName(name string) string {
	return fmt.Sprintf("%s%s", RpcConsumer, name)
}