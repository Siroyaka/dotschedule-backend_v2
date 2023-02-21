package utility

import (
	"encoding/json"
	"io"
	"reflect"
	"regexp"
	"strings"
)

const (
	EmptyString = ""

	replaceConst        = "@@"
	config_timeFormat   = "TIME_FORMAT"
	config_dateFormat   = "DATE_FORMAT"
	config_monthFormat  = "MONTH_FORMAT"
	config_timeLocation = "LOCAL_LOCATION"
)

func ReplaceConstString(src, repl, target string) string {
	target = replaceConst + target + replaceConst
	return strings.ReplaceAll(src, target, repl)
}

func ReplaceString(src, repl, regex string) string {
	rel := regexp.MustCompile(regex)
	return rel.ReplaceAllString(src, repl)
}

func Contains[T comparable](list []T, target T) bool {
	for _, d := range list {
		if d == target {
			return true
		}
	}
	return false
}

func fromArray[X comparable](data []X) (res []interface{}) {
	for _, s := range data {
		res = append(res, ToInterfaceSlice(s)...)
	}
	return res
}

func ToInterfaceSlice(l ...interface{}) (list []interface{}) {
	for _, v := range l {
		switch value := v.(type) {
		case []string:
			list = append(list, fromArray(value)...)
		case []int:
			list = append(list, fromArray(value)...)
		case []bool:
			list = append(list, fromArray(value)...)
		case []interface{}:
			list = append(list, ToInterfaceSlice(value...)...)
		default:
			list = append(list, value)
		}
	}

	return list
}

func JsonUnmarshal[X any](jsonData string) (result X, ierr IError) {
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		ierr = NewError(err.Error(), ERR_JSONPARSE)
	}
	return
}

func JsonDecode[X any](reader io.Reader) (result X, ierr IError) {
	if err := json.NewDecoder(reader).Decode(&result); err != nil {
		ierr = NewError(err.Error(), ERR_JSONPARSE)
	}
	return
}

func ConvertFromInterfaceType[X any](value interface{}) (res X, err IError) {
	switch ci := value.(type) {
	case X:
		res = ci
		err = nil
		return
	default:
		err = NewError("", ERR_TYPEPARSE, reflect.TypeOf(value))
		return
	}
}

type HashSet[X comparable] struct {
	set map[X]struct{}
}

func NewHashSet[X comparable]() HashSet[X] {
	set := make(map[X]struct{})
	return HashSet[X]{set: set}
}

func (hash *HashSet[X]) Set(value X) bool {
	if _, ok := hash.set[value]; ok {
		return false
	}
	hash.set[value] = struct{}{}
	return true
}

func (hash HashSet[X]) Has(value X) bool {
	_, ok := hash.set[value]
	return ok
}

func (hash HashSet[X]) List() []X {
	var list []X
	for k := range hash.set {
		list = append(list, k)
	}
	return list
}

func (hash HashSet[X]) Any() bool {
	for range hash.set {
		return true
	}
	return false
}

func ConditionalOperator[X any](flg bool, trueValue, falseValue X) X {
	if flg {
		return trueValue
	} else {
		return falseValue
	}
}
