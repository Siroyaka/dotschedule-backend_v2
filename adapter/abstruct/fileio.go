package abstruct

import (
	"io"

	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type FileInfo interface {
	ReadFile(string) (string, utility.IError)
	ReadFileLine(string) ([]string, utility.IError)
	Any(string) bool
	FileList(string) ([]string, utility.IError)
	IsDirectory(string) bool
}

type FileReader[X any] interface {
	Read(string, func(io.Reader) (X, utility.IError)) (X, utility.IError)
	FileList(string) ([]string, utility.IError)
}
