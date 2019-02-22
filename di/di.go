package di

import (
	"fmt"
	"reflect"
	"strings"
	"tianwei.pro/micro/rpc/client"
	"tianwei.pro/micro/rpc/server"
)

type Dependency interface {
	Open() error
	Close()
}

// Container contains dependencies by name.
type Container struct {
	allowCover   bool
	dependencies map[string]interface{}
}

// New creates new Container instance.
func New() *Container {
	return &Container{
		allowCover:   false,
		dependencies: make(map[string]interface{}),
	}
}

// Provide registers a dependency.
func (c *Container) Provide(name string, dependency interface{}) {
	if _, ok := c.dependencies[name]; ok && !c.allowCover {
		panic(fmt.Sprintf("provider not allow cover. name: %s", name))
	}
	c.dependencies[name] = dependency
}

// GetByName returns a dependency by name.
func (c *Container) GetByName(name string) interface{} {
	return c.dependencies[name]
}

// GetByType returns a dependency by type.
func (c *Container) GetByType(typ interface{}) interface{} {
	t := reflect.TypeOf(typ)
	for _, d := range c.dependencies {
		dt := reflect.TypeOf(d)
		if t.AssignableTo(dt) {
			return d
		}
	}
	return nil
}

// Resolve decorates objects with dependencies and initializes them.
func (c *Container) Resolve() error {
	for key, d := range c.dependencies {
		if strings.Contains(key, RpcConsumer) {
			client.UserService(d)
		}
	}
	for _, d := range c.dependencies {
		c.inject(d)
	}
	for _, d := range c.dependencies {
		if dep, ok := d.(Dependency); ok {
			err := dep.Open()
			if err != nil {
				panic(err)
			}
		}
	}
	return nil
}

func (c *Container) ResolverRpcServer() {
	for key, d := range c.dependencies {
		if strings.Contains(key, RpcProvider) {
			server.AddInstanceMethods(d)
		}
	}
}

// Close closes all dependencies.
func (c *Container) Close() {
	for _, d := range c.dependencies {
		if dep, ok := d.(Dependency); ok {
			dep.Close()
		}
	}
}

func (c *Container) inject(obj interface{}) {
	t := reflect.TypeOf(obj).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		inject := field.Tag.Get("inject")
		if inject == "" {
			continue
		}
		required := true
		if field.Tag.Get("required") == "false" {
			required = false
		}

		dependency := c.GetByName(inject)
		if dependency != nil {
			reflect.ValueOf(obj).Elem().Field(i).Set(reflect.ValueOf(dependency))
		} else if required {
			panic(fmt.Sprintf("Colud not inject field: %s", inject))
		}
	}
}
