package utility_test

import (
	"testing"

	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

func TestToInterfaceSliceBasic(t *testing.T) {
	a := utility.ToInterfaceSlice("a", 1, false)
	if a[0] != "a" {
		t.Error()
	}
	if a[1] != 1 {
		t.Error()
	}
	if a[2] != false {
		t.Error("Array 3 is wrong value", a[2])
	}
}

func TestToInterfaceSliceArray(t *testing.T) {
	stringSlice := []string{"value", "value2", "value3"}
	a := utility.ToInterfaceSlice("a", stringSlice)
	if a[0] != "a" {
		t.Error()
	}
	for i := 0; i < len(stringSlice); i++ {
		if a[i+1] != stringSlice[i] {
			t.Error("array wrong", a[i+1], stringSlice[i])
		}
	}
}
