package utility

type IConfig interface {
	Read(string) string
	ReadInteger(string) int
	ReadChild(string) IConfig
	ReadStringList(string) []string
	Has(string) bool
	ConfigList() []string
}
