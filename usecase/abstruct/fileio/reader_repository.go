package fileio

import "github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"

type ReaderRepository[X any] interface {
	FileList(string) ([]string, utilerror.IError)
	ReadJson(string) (X, utilerror.IError)
}
