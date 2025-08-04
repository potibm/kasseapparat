package errors

type BasicError struct {
	Code     int
	Message  string
	Detail   string
	CauseErr error
}

func (e *BasicError) Error() string   { return e.Message }
func (e *BasicError) Details() string { return e.Detail }
func (e *BasicError) StatusCode() int { return e.Code }
func (e *BasicError) Cause() error    { return e.CauseErr }
func (e *BasicError) WithCause(err error) *BasicError {
	return &BasicError{Code: e.Code, Message: e.Message, Detail: e.Detail, CauseErr: err}
}
func (e *BasicError) WithCauseMsg(err error) *BasicError {
	return &BasicError{Code: e.Code, Message: e.Message, Detail: err.Error(), CauseErr: err}
}
func (e *BasicError) WithMsg(message string) *BasicError {
	return &BasicError{Code: e.Code, Message: e.Message, Detail: message, CauseErr: e.CauseErr}
}

func (e *BasicError) Unwrap() error { return e.CauseErr }
