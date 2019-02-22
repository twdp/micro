package yaml

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/beego/goyaml2"
	yaml2 "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type Config struct{}

// ConfigContainer A Config represents the yaml configuration.
type ConfigContainer struct {
	data map[string]interface{}
	sync.RWMutex
}

func (yaml *Config) Parse(filename string) (y config.Configer, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	res := make(map[string]interface{})
	if err := yaml2.Unmarshal(data, &res); err != nil {
		panic(err)
	}

	y = &ConfigContainer{
		data: res,
	}
	return
}

func (yaml *Config) ParseData(data []byte) (config.Configer, error) {
	res := make(map[string]interface{})
	if err := yaml2.Unmarshal(data, &res); err != nil {
		panic(err)
	}

	return &ConfigContainer{
		data: res,
	}, nil
}

// Bool returns the boolean value for a given key.
func (c *ConfigContainer) Bool(key string) (bool, error) {
	v, err := c.getData(key)
	if err != nil {
		return false, err
	}
	return config.ParseBool(v)
}

// DefaultBool return the bool value if has no error
// otherwise return the defaultval
func (c *ConfigContainer) DefaultBool(key string, defaultval bool) bool {
	v, err := c.Bool(key)
	if err != nil {
		return defaultval
	}
	return v
}

// Int returns the integer value for a given key.
func (c *ConfigContainer) Int(key string) (int, error) {
	if v, err := c.getData(key); err != nil {
		return 0, err
	} else if vv, ok := v.(int); ok {
		return vv, nil
	} else if vv, ok := v.(int64); ok {
		return int(vv), nil
	}
	return 0, errors.New("not int value")
}

// DefaultInt returns the integer value for a given key.
// if err != nil return defaultval
func (c *ConfigContainer) DefaultInt(key string, defaultval int) int {
	v, err := c.Int(key)
	if err != nil {
		return defaultval
	}
	return v
}

// Int64 returns the int64 value for a given key.
func (c *ConfigContainer) Int64(key string) (int64, error) {
	if v, err := c.getData(key); err != nil {
		return 0, err
	} else if vv, ok := v.(int64); ok {
		return vv, nil
	}
	return 0, errors.New("not bool value")
}

// DefaultInt64 returns the int64 value for a given key.
// if err != nil return defaultval
func (c *ConfigContainer) DefaultInt64(key string, defaultval int64) int64 {
	v, err := c.Int64(key)
	if err != nil {
		return defaultval
	}
	return v
}

// Float returns the float value for a given key.
func (c *ConfigContainer) Float(key string) (float64, error) {
	if v, err := c.getData(key); err != nil {
		return 0.0, err
	} else if vv, ok := v.(float64); ok {
		return vv, nil
	} else if vv, ok := v.(int); ok {
		return float64(vv), nil
	} else if vv, ok := v.(int64); ok {
		return float64(vv), nil
	}
	return 0.0, errors.New("not float64 value")
}

// DefaultFloat returns the float64 value for a given key.
// if err != nil return defaultval
func (c *ConfigContainer) DefaultFloat(key string, defaultval float64) float64 {
	v, err := c.Float(key)
	if err != nil {
		return defaultval
	}
	return v
}

// String returns the string value for a given key.
func (c *ConfigContainer) String(key string) string {
	if v, err := c.getData(key); err == nil {
		if vv, ok := v.(string); ok {
			return vv
		}
	}
	return ""
}

// DefaultString returns the string value for a given key.
// if err != nil return defaultval
func (c *ConfigContainer) DefaultString(key string, defaultval string) string {
	v := c.String(key)
	if v == "" {
		return defaultval
	}
	return v
}

// Strings returns the []string value for a given key.
func (c *ConfigContainer) Strings(key string) []string {
	v := c.String(key)
	if v == "" {
		return nil
	}
	return strings.Split(v, ";")
}

// DefaultStrings returns the []string value for a given key.
// if err != nil return defaultval
func (c *ConfigContainer) DefaultStrings(key string, defaultval []string) []string {
	v := c.Strings(key)
	if v == nil {
		return defaultval
	}
	return v
}

// GetSection returns map for the given section
func (c *ConfigContainer) GetSection(section string) (map[string]string, error) {

	if v, ok := c.data[section]; ok {
		return v.(map[string]string), nil
	}
	return nil, errors.New("not exist section")
}

// SaveConfigFile save the config into file
func (c *ConfigContainer) SaveConfigFile(filename string) (err error) {
	// Write configuration file by filename.
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	err = goyaml2.Write(f, c.data)
	return err
}

// Set writes a new value for key.
func (c *ConfigContainer) Set(key, val string) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = val
	return nil
}

// DIY returns the raw value by a given key.
func (c *ConfigContainer) DIY(key string) (v interface{}, err error) {
	return c.getData(key)
}

func (c *ConfigContainer) getData(key string) (interface{}, error) {

	if len(key) == 0 {
		return nil, errors.New("key is empty")
	}
	c.RLock()
	defer c.RUnlock()

	keys := strings.Split(key, ".")
	tmpData := c.data
	for idx, k := range keys {
		if v, ok := tmpData[k]; ok {
			switch v.(type) {
			case map[string]interface{}:
				{
					tmpData = v.(map[string]interface{})
					if idx == len(keys)-1 {
						return tmpData, nil
					}
				}
			default:
				{
					return v, nil
				}

			}
		}
	}
	return nil, fmt.Errorf("not exist key %q", key)
}

func init() {
	config.Register("yaml2", &Config{})
}
