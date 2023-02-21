package fileio

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type DirectoryFileList struct {
	reader abstruct.FileReader[reference.VoidStruct]
}

func NewDirectoryFileList(reader abstruct.FileReader[reference.VoidStruct]) DirectoryFileList {
	return DirectoryFileList{
		reader: reader,
	}
}

func (repos DirectoryFileList) Execute(dirPath string) ([]string, utility.IError) {
	nameList, err := repos.reader.FileList(dirPath)
	if err != nil {
		return []string{}, err.WrapError()
	}
	return nameList, nil
}
