package config

type IConfig interface {
	Read(string) string
	ReadInteger(string) int
	ReadBoolean(string) bool
	ReadChild(string) IConfig
	ReadStringList(string) []string
	Has(string) bool
	ConfigList() []string
}

type ConfigData map[string]interface{}
