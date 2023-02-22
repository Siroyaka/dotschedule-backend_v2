package infrastructure

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"

	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type FileReader[X any] struct {
}

func NewFileReader[X any]() FileReader[X] {
	return FileReader[X]{}
}

func (fr FileReader[X]) Read(filePath string, f func(io.Reader) (X, utility.IError)) (data X, ierr utility.IError) {
	osFile, err := os.Open(filePath)
	defer osFile.Close()
	if err != nil {
		ierr = utility.NewError(err.Error(), utility.ERR_FILE_READ, filePath)
		return
	}
	data, ierr = f(osFile)
	if ierr != nil {
		ierr = ierr.WrapError()
		return
	}

	return
}

func (fr FileReader[X]) FileList(dirPath string) (list []string, ierr utility.IError) {
	f, err := os.Stat(dirPath)
	if err != nil {
		ierr = utility.NewError(err.Error(), utility.ERR_DIRECTORY_READ, dirPath)
		return
	}
	if !f.IsDir() {
		ierr = utility.NewError("there is not directory", utility.ERR_DIRECTORY_READ, dirPath)
		return
	}
	items, err := ioutil.ReadDir(dirPath)
	if err != nil {
		ierr = utility.NewError(err.Error(), utility.ERR_DIRECTORY_READ, dirPath)
		return
	}
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		list = append(list, item.Name())
	}
	return
}

type FileInfo struct {
}

func NewFileInfo() FileInfo {
	return FileInfo{}
}

func (info FileInfo) IsDirectory(path string) bool {
	f, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return f.IsDir()
}

func (io FileInfo) Any(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (io FileInfo) FileList(dirPath string) (list []string, ierr utility.IError) {
	f, err := os.Stat(dirPath)
	if err != nil {
		ierr = utility.NewError(err.Error(), utility.ERR_DIRECTORY_READ, dirPath)
		return
	}
	if !f.IsDir() {
		ierr = utility.NewError("there is not directory", utility.ERR_DIRECTORY_READ, dirPath)
		return
	}
	items, err := ioutil.ReadDir(dirPath)
	if err != nil {
		ierr = utility.NewError(err.Error(), utility.ERR_DIRECTORY_READ, dirPath)
		return
	}
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		list = append(list, item.Name())
	}
	return
}

func (io FileInfo) ReadFile(filePath string) (string, utility.IError) {
	osFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", utility.NewError(err.Error(), utility.ERR_FILE_READ, filePath)
	}

	return string(osFile), nil
}

func (io FileInfo) ReadFileLine(filePath string) ([]string, utility.IError) {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return []string{}, utility.NewError(err.Error(), utility.ERR_FILE_READ, filePath)
	}
	res := []string{}
	fr := bufio.NewScanner(f)
	err = fr.Err()
	if err != nil {
		return []string{}, utility.NewError(err.Error(), utility.ERR_FILE_READ, filePath)
	}
	for fr.Scan() {
		res = append(res, fr.Text())
	}
	return res, nil
}
