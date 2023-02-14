package abstruct

import "github.com/Siroyaka/dotschedule-backend_v2/utility"

type HTTPRequest interface {
	Get(string) (HTTPResponse, utility.IError)
	Post(HTTPPostParams) (HTTPResponse, utility.IError)
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
