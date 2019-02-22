package conf

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"testing"
	_ "tianwei.pro/micro/conf/yaml"
)

func TestDiy(t *testing.T) {
	c, err := config.NewConfig("yaml", "conf.yaml")
	if err != nil {
		logs.Error("parse config failed. %v", err)
		panic(err)
	}
	fmt.Println(c.String("client.address.product"))
	addresss, _ := c.DIY("client.address")
	for k, v := range addresss.(map[string]interface{}) {
		fmt.Println(k)
		fmt.Println(v)
	}
}
