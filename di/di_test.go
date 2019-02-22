package di

import (
	"fmt"
	"github.com/hprose/hprose-golang/rpc"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type Database interface {
	GetValue() int
}

type databaseImpl struct {
	//Started bool
}

func (d *databaseImpl) GetValue() int {
	return 42
}

func (d *databaseImpl) Open() error {
	//d.Started = true
	return nil
}

func (d *databaseImpl) Close() {
}

type service struct {
	Database Database `inject:"db"`
	Started  bool
}

func (s *service) Open() error {
	s.Started = true
	return nil
}

func (s *service) Close() {
}

func TestDependencyInjection(t *testing.T) {
	c := New()
	c.Provide("db", &databaseImpl{})
	c.Provide("s", &service{})

	assert.NoError(t, c.Resolve())

	 //c.GetByName("db").(*databaseImpl)
	//assert.True(t, db.Started)
	s, _ := c.GetByName("s").(*service)
	assert.True(t, s.Started)
	assert.Equal(t, 42, s.Database.GetValue())
}

type A struct {
	B *B
}

func (a *A) Value() int {
	return a.B.Value()
}

type B struct {
}

func (b *B) Value() int {
	return 42
}

func TestWithStruct(t *testing.T) {
	c := New()
	c.Provide("a", &A{})
	c.Provide("b", &B{})
	assert.NoError(t, c.Resolve())

	a := c.GetByType(&A{}).(*A)
	assert.Equal(t, 42, a.Value())
}

type service2 struct {
	Database *data `inject:"db"`
	Started  bool
}

type data struct {
	Hello func(a string) (string, error)
}

func hello(a string) string{
	return "aaaaaa"
}
var mapping = make(map[string]interface{})
func TestContainer_Close(t *testing.T) {
	mapping["data"] = &data{}

	server := rpc.NewTCPServer("")
	server.AddFunction("hello", hello)
	server.Handle()
	defer server.Close()
	client := rpc.NewClient(server.URI())
	defer client.Close()

	s := &service2{}
	tt := reflect.TypeOf(s).Elem()
	for i := 0; i < tt.NumField(); i++ {
		if i == 0 {
			f := reflect.ValueOf(s).Elem().Field(i).Type().Elem()
			fmt.Println(f.Name())
			var ss *databaseImpl
			fmt.Println(reflect.ValueOf(&ss))
			v := reflect.ValueOf(&f)
			fmt.Println(v.Kind() == reflect.Ptr)
			var mappii = mapping[f.Name()]
			client.UseService(mappii)
			fmt.Println(mappii)
			fmt.Println((mappii.(*data)).Hello("aa"))
		}
	}

	fmt.Println(s)
}

func buildRemoteService(v reflect.Value) {
	v = v.Elem()
	t := v.Type()
	et := t
	if et.Kind() == reflect.Ptr {
		et = et.Elem()
	}
	ptr := reflect.New(et)
	obj := ptr.Elem()

	if t.Kind() == reflect.Ptr {
		v.Set(ptr)
	} else {
		v.Set(obj)
	}
}

