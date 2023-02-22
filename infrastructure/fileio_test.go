package infrastructure_test

import (
	"testing"

	"github.com/Siroyaka/dotschedule-backend_v2/infrastructure"
)

func TestFileList(t *testing.T) {
	io := infrastructure.NewFileInfo()
	list, err := io.FileList("./")
	if err != nil {
		t.Error(err)
	}
	for _, v := range list {
		t.Log(v)
	}
}
