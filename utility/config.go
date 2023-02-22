package utility

import (
	"github.com/Siroyaka/dotschedule-backend_v2/utility/config"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
)

type ConfigData map[string]interface{}

type Config struct {
	configData ConfigData
}

func NewConfig(data ConfigData) config.IConfig {
	if data == nil {
		return Config{make(ConfigData)}
	}
	return Config{configData: data}
}

func (c Config) Has(key string) bool {
	_, has := c.configData[key]
	return has
}

func (c Config) check(k string) {
	if !c.Has(k) {
		logger.Fatal(utilerror.New("", utilerror.ERR_CONFIG_NOTFOUND, k))
		panic(utilerror.New("", utilerror.ERR_CONFIG_NOTFOUND, k))
	}
}

func configRead[X any](data interface{}) (result X, err utilerror.IError) {
	switch value := data.(type) {
	case X:
		result = value
		err = nil
		return
	default:
		err = utilerror.New("config type error", "")
		return
	}
}

func configListRead[X any](list []interface{}) (result []X, err utilerror.IError) {
	for _, v := range list {
		r, e := configRead[X](v)
		if e != nil {
			err = e
			return
		}
		result = append(result, r)
	}
	err = nil
	return
}

func (c Config) Read(k string) string {
	c.check(k)
	result, err := configRead[string](c.configData[k])
	if err != nil {
		panic(err)
	}
	return result
}

func (c Config) ReadInteger(k string) int {
	c.check(k)
	result, err := configRead[float64](c.configData[k])
	if err != nil {
		panic(err)
	}
	return int(result)
}

func (c Config) ReadBoolean(k string) bool {
	c.check(k)
	result, err := configRead[bool](c.configData[k])
	if err != nil {
		panic(err)
	}
	return result
}

func (c Config) ReadChild(k string) config.IConfig {
	c.check(k)
	result, err := configRead[map[string]interface{}](c.configData[k])
	if err != nil {
		panic(err)
	}
	return Config{
		configData: result,
	}
}

func (c Config) ReadStringList(k string) []string {
	c.check(k)
	result, err := configRead[[]interface{}](c.configData[k])
	if err != nil {
		panic(err)
	}
	list, err := configListRead[string](result)
	if err != nil {
		panic(err)
	}
	return list
}

func (c Config) ConfigList() []string {
	keys := make([]string, 0, len(c.configData))
	for k := range c.configData {
		keys = append(keys, k)
	}
	return keys
}
