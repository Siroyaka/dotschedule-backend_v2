package fileio

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/fileio"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
)

type ReaderRepository[X any] struct {
	reader abstruct.FileReader[X]
}

func NewReaderRepository[X any](reader abstruct.FileReader[X]) fileio.ReaderRepository[X] {
	return ReaderRepository[X]{
		reader: reader,
	}
}

func (repos ReaderRepository[X]) FileList(dirPath string) ([]string, utilerror.IError) {
	nameList, err := repos.reader.FileList(dirPath)
	if err != nil {
		return []string{}, err.WrapError()
	}
	return nameList, nil
}

func (repos ReaderRepository[X]) ReadJson(filePath string) (result X, err utilerror.IError) {
	data, err := repos.reader.Read(filePath, utility.JsonDecode[X])
	if err != nil {
		err = err.WrapError()
		return
	}
	return data, nil
}
