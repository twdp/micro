// +build !rpc

package conf

import "github.com/astaxie/beego"

func init() {
	Conf = beego.AppConfig
}