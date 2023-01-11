package fileio

import (
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type ReaderRepository[X any] interface {
	FileList(string) ([]string, utility.IError)
	ReadJson(string) (X, utility.IError)
}
