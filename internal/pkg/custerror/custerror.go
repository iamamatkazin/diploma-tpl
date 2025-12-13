package custerror

type Error struct {
	Message string
	Code    int
}

func (e *Error) Error() string {
	return e.Message
}

func New(code int, message string) error {
	return &Error{Code: code, Message: message}
}
