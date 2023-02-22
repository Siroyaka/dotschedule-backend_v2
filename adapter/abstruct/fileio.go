package abstruct

import (
	"io"

	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
)

type FileInfo interface {
	ReadFile(string) (string, utilerror.IError)
	ReadFileLine(string) ([]string, utilerror.IError)
	Any(string) bool
	FileList(string) ([]string, utilerror.IError)
	IsDirectory(string) bool
}

type FileReader[X any] interface {
	Read(string, func(io.Reader) (X, utilerror.IError)) (X, utilerror.IError)
	FileList(string) ([]string, utilerror.IError)
}
