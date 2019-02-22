package single

import "tianwei.pro/micro/di"

var container *di.Container

func init() {
	container = di.New()
}

func Provide(name string, dependency interface{}) {
	container.Provide(name, dependency)
}

func Resolve() {
	if err := container.Resolve(); err != nil {
		panic(err)
	}
}

func GetByName(name string) interface{} {
	return container.GetByName(name)
}