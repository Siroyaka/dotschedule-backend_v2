package repository

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type IConfigLoader interface {
	ReadJsonConfig(string) (utility.IConfig, utility.IError)
}

type ConfigLoader struct {
	fileReader abstruct.FileReader[utility.ConfigData]
}

func NewConfigLoader(io abstruct.FileReader[utility.ConfigData]) IConfigLoader {
	return ConfigLoader{fileReader: io}
}

func (loader ConfigLoader) ReadJsonConfig(filePath string) (utility.IConfig, utility.IError) {
	configData, err := loader.fileReader.Read(filePath, utility.JsonDecode[utility.ConfigData])
	if err != nil {
		return nil, err.WrapError()
	}
	return utility.NewConfig(configData), nil
}