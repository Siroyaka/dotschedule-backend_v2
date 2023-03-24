package utilerror

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type Error struct {
	errorCode    string
	errorType    ErrorType
	errorMessage string
	bottomError  error
	callStacks   []callStack
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
	case ERR_SQL_DATAUPDATE_COUNT0:
		code = err_sql_dataupdate_count0_code
		err = errors.New(err_sql_dataupdate_count0_msg)
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

func New(message string, errType ErrorType, errInfo ...any) IError {
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
	error := Error{
		errorCode:    code,
		errorType:    errType,
		errorMessage: issue,
		bottomError:  err,
		callStacks:   []callStack{newStack},
	}

	return error
}

type callStack struct {
	function string
	message  string
}

func (cs callStack) show() string {
	return fmt.Sprintf("  %s: %s", cs.function, cs.message)
}
