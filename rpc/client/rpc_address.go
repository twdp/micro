// +build rpc

package client

func initAddress() {
	address = conf.Conf.Strings("client.address")
	if len(address) == 0 {
		panic("client address not found")
	}
}