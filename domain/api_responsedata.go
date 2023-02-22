package domain

import (
	"encoding/json"

	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
)

type APIResponseData[X any] struct {
	status        string
	length        int
	error_message string
	response_data X
}

func NewAPIResponseData[X any](status string, length int, errorMessage string, responseData X) APIResponseData[X] {
	return APIResponseData[X]{
		status:        status,
		length:        length,
		error_message: errorMessage,
		response_data: responseData,
	}
}

func (ad APIResponseData[X]) ToJson() (string, utilerror.IError) {
	d := make(map[string]interface{})
	d["status"] = ad.status
	d["length"] = ad.length
	d["error_message"] = ad.error_message
	d["response_data"] = ad.response_data

	responseJson, err := json.Marshal(d)
	if err != nil {
		return "", utilerror.New(err.Error(), utilerror.ERR_JSONPARSE)
	}
	return string(responseJson), nil
}
