package utilerror

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
