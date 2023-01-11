package utility

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type ErrorType string

const (
	// 0****
	ERR_UNKNOWN      ErrorType = "UNKNOWN ERROR"
	err_unknown_code           = "00000"
	err_unknown_msg            = "unknown error"

	ERR_CONFIG_NOTFOUND      ErrorType = "CONFIG_NOTFOUND"
	err_config_notfound_code           = "00001"
	err_config_notfound_msg            = "config not found"

	// 1**** domain or utility error
	ERR_LOAD_TIMELOCATION      ErrorType = "LOAD_TIMELOCATION_ERROR"
	err_load_timelocation_code           = "10000"
	err_load_timelocation_msg            = "time location loading error"

	ERR_TIME_PARSE      ErrorType = "TIME_PARSE_ERROR"
	err_time_parse_code           = "10001"
	err_time_parse_msg            = "text to time parse error"

	ERR_RSS_PARSE      ErrorType = "RSS_PARSE_ERROR"
	err_rss_parse_code           = "10100"
	err_rss_parse_msg            = "rss data parse error"

	// 8**** interface error
	ERR_SQL_PREPARE      ErrorType = "SQL_PREPARE_ERROR"
	err_sql_prepare_code           = "80001"
	err_sql_prepare_msg            = "sql prepare error"

	ERR_SQL_QUERY      ErrorType = "SQL_QUERY_ERROR"
	err_sql_query_code           = "80002"
	err_sql_query_msg            = "sql query error"

	ERR_SQL_DATASCAN      ErrorType = "SQL_DATASCAN_ERROR"
	err_sql_datascan_code           = "80003"
	err_sql_datascan_msg            = "sql datascan error"

	ERR_HTTP_REQUEST_ERROR      ErrorType = "HTTP_REQUEST_ERROR"
	err_http_request_error_code           = "81000"
	err_http_request_error_msg            = "http request error"

	ERR_HTTP_REQUEST_TIMEOUT      ErrorType = "HTTP_REQUEST_TIMEOUT"
	err_http_request_timeout_code           = "81001"
	err_http_request_timeout_msg            = "http request timeout"

	ERR_HTTP_BODY_READERROR              ErrorType = "HTTP_BODY_READERROR"
	err_http_request_body_readerror_code           = "81101"
	err_http_request_body_readerror_msg            = "http body read error"

	ERR_FILE_READ      ErrorType = "FILE_READ_ERROR"
	err_file_read_code           = "80100"
	err_file_read_msg            = "file read error"

	ERR_DIRECTORY_READ      ErrorType = "DIRECTORY_READ_ERROR"
	err_directory_read_code           = "80101"
	err_directory_read_msg            = "directory read error"

	// 9**** common error
	ERR_OUTOFINDEX      ErrorType = "OUT_OF_INDEX"
	err_outofindex_code           = "90000"
	err_outofindex_msg            = "out of index"

	ERR_INVALIDVALUE      ErrorType = "INVALID_VALUE"
	err_invalidvalue_code           = "90001"
	err_invalidvalue_msg            = "invalid value"

	ERR_JSONPARSE      ErrorType = "JSON_PARSE"
	err_jsonparse_code           = "91000"
	err_jsonparse_msg            = "json parse error"

	ERR_TYPEPARSE      ErrorType = "TYPE_PARSE"
	err_typeparse_code           = "92000"
	err_typeparse_msg            = "value type parse error"
)

type Error struct {
	errorCode    string
	errorType    ErrorType
	errorMessage string
	bottomError  error
	callStacks   []callStack
}

type IError interface {
	Error() string
	WrapError(...string) IError
	TypeIs(ErrorType) bool
}

func (er Error) TypeIs(errType ErrorType) bool {
	return er.errorType == errType
}

func (er Error) Error() string {
	var stackMessages []string
	for _, cs := range er.callStacks {
		stackMessages = append(stackMessages, cs.show())
	}
	return fmt.Sprintf("%s:%s %s\n%s", er.errorCode, er.bottomError, er.errorMessage, strings.Join(stackMessages, "\n"))
}

func (er Error) WrapError(msgs ...string) IError {
	functionName := "??"
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		functionInfo := runtime.FuncForPC(pc)
		functionName = functionInfo.Name()
	}
	newStack := callStack{
		function: functionName,
		message:  strings.Join(msgs, " "),
	}
	callStacks := append(er.callStacks, newStack)
	return Error{
		errorCode:    er.errorCode,
		errorType:    er.errorType,
		errorMessage: er.errorMessage,
		bottomError:  er.bottomError,
		callStacks:   callStacks,
	}
}

func createError(errType ErrorType, errorInfo ...any) (code string, issue string, err error) {
	switch errType {
	case ERR_CONFIG_NOTFOUND:
		code = err_config_notfound_code
		err = errors.New(err_config_notfound_msg)
		if len(errorInfo) == 1 {
			issue = fmt.Sprintf("key:%s", errorInfo[0])
		}
		return

	case ERR_LOAD_TIMELOCATION:
		code = err_config_notfound_code
		err = errors.New(err_load_timelocation_msg)
		if len(errorInfo) == 1 {
			issue = fmt.Sprintf("location:%s", errorInfo[0])
		}
		return
	case ERR_TIME_PARSE:
		code = err_time_parse_code
		err = errors.New(err_time_parse_msg)
		if len(errorInfo) == 3 {
			issue = fmt.Sprintf("str:%s location:%s format:%s", errorInfo[0], errorInfo[1], errorInfo[2])
		}
		return
	case ERR_SQL_PREPARE:
		code = err_sql_prepare_code
		err = errors.New(err_sql_prepare_msg)
		if len(errorInfo) == 1 {
			issue = fmt.Sprintf("query:%s", errorInfo[0])
		}
		return
	case ERR_SQL_QUERY:
		code = err_sql_query_code
		err = errors.New(err_sql_query_msg)
		if len(errorInfo) == 1 {
			issue = fmt.Sprintf("query:%s", errorInfo[0])
		}
		return
	case ERR_SQL_DATASCAN:
		code = err_sql_datascan_code
		err = errors.New(err_sql_datascan_msg)
		return
	case ERR_OUTOFINDEX:
		code = err_outofindex_code
		err = errors.New(err_outofindex_msg)
		if len(errorInfo) == 3 {
			issue = fmt.Sprintf("%s len is %s. out of index. index value:%s", errorInfo[0], errorInfo[1], errorInfo[2])
		}
		return
	case ERR_INVALIDVALUE:
		code = err_invalidvalue_code
		err = errors.New(err_invalidvalue_msg)
		if len(errorInfo) == 2 {
			issue = fmt.Sprintf("%s is invalid value: %s", errorInfo[1], errorInfo[0])
		}
		return
	case ERR_JSONPARSE:
		code = err_jsonparse_code
		err = errors.New(err_jsonparse_msg)
		return
	case ERR_TYPEPARSE:
		code = err_typeparse_code
		err = errors.New(err_typeparse_msg)
		if len(errorInfo) == 1 {
			issue = fmt.Sprintf("value type: %s", errorInfo[0])
		}
		return
	case ERR_FILE_READ:
		code = err_file_read_code
		err = errors.New(err_file_read_msg)
		if len(errorInfo) == 1 {
			issue = fmt.Sprintf("%s read error", errorInfo[0])
		}
		return
	case ERR_DIRECTORY_READ:
		code = err_directory_read_code
		err = errors.New(err_directory_read_msg)
		if len(errorInfo) == 1 {
			issue = fmt.Sprintf("%s read error", errorInfo[0])
		}
		return
	case ERR_HTTP_REQUEST_ERROR:
		code = err_http_request_error_code
		err = errors.New(err_http_request_error_msg)
		return
	case ERR_HTTP_REQUEST_TIMEOUT:
		code = err_http_request_timeout_code
		err = errors.New(err_http_request_timeout_msg)
		return
	case ERR_HTTP_BODY_READERROR:
		code = err_http_request_body_readerror_code
		err = errors.New(err_http_request_body_readerror_msg)
		return
	case ERR_RSS_PARSE:
		code = err_rss_parse_code
		err = errors.New(err_rss_parse_msg)
		return
	default:
		code = err_unknown_code
		err = errors.New(err_unknown_msg)
		return
	}
}

func NewError(message string, errType ErrorType, errInfo ...any) IError {
	functionName := "??"
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		functionInfo := runtime.FuncForPC(pc)
		functionName = functionInfo.Name()
	}
	code, issue, err := createError(errType, errInfo...)
	newStack := callStack{
		function: functionName,
		message:  message,
	}
	a := Error{
		errorCode:    code,
		errorType:    errType,
		errorMessage: issue,
		bottomError:  err,
		callStacks:   []callStack{newStack},
	}

	return a
}

type callStack struct {
	function string
	message  string
}

func (cs callStack) show() string {
	return fmt.Sprintf("  %s: %s", cs.function, cs.message)
}
