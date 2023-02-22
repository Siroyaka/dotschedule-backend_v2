package utilerror

type IError interface {
	Error() string
	WrapError(...string) IError
	TypeIs(ErrorType) bool
}

type ErrorType string
