package abstruct

import "github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"

type HTTPRequest interface {
	Get(string) (HTTPResponse, utilerror.IError)
	Post(HTTPPostParams) (HTTPResponse, utilerror.IError)
	SetTimeout(int)
}

type HTTPResponse interface {
	Status() string
	StatusCode() int
	Body() string
}

type HTTPPostParams interface {
	ContentType() string
	Url() string
	Content() string
}
