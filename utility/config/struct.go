package config

import "log"

var (
	config      IConfig
	projectName string
)

func Setup(n string, c IConfig) {
	projectName = n
	config = c
}

func ConfigIsNil() bool {
	if config == nil {
		log.Printf("ERROR\t%s\n", "Config is Nothing. Need Setting Config.")
		return true
	}
	return false
}

func Has(key string) bool {
	if ConfigIsNil() {
		return false
	}
	return config.Has(key)
}

func Read(k string) string {
	if ConfigIsNil() {
		return ""
	}
	return config.Read(k)
}

func ReadInteger(k string) int {
	if ConfigIsNil() {
		return 0
	}
	return config.ReadInteger(k)
}

func ReadChild(k string) IConfig {
	if ConfigIsNil() {
		return nil
	}
	return config.ReadChild(k)
}

func ReadProjectConfig() IConfig {
	if ConfigIsNil() {
		return nil
	}
	return config.ReadChild(projectName)
}

func ReadStringList(k string) []string {
	if ConfigIsNil() {
		return []string{}
	}
	return config.ReadStringList(k)
}

func ConfigList() []string {
	if ConfigIsNil() {
		return []string{}
	}
	return config.ConfigList()
}
