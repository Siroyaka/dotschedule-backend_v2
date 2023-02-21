package fileio

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type JsonFileReader[X any] struct {
	reader abstruct.FileReader[X]
}

func NewJsonFileReader[X any](reader abstruct.FileReader[X]) JsonFileReader[X] {
	return JsonFileReader[X]{
		reader: reader,
	}
}

func (repos JsonFileReader[X]) Execute(filePath string) (result X, err utility.IError) {
	data, err := repos.reader.Read(filePath, utility.JsonDecode[X])
	if err != nil {
		err = err.WrapError()
		return
	}
	return data, nil
}
