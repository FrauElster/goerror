package goerror

import (
	"encoding/json"
	"fmt"
	"runtime"
)

type TraceableError interface {
	error
	Unwrap() error
	Is(error) bool
	MarshalJSON() ([]byte, error)
	WithMessage(string) TraceableError
	WithError(error) TraceableError
	WithOrigin() TraceableError
	GetMessage() string
}

func New(code, message string) traceableErrorImpl {
	return traceableErrorImpl{Code: code, Message: message}
}

var _ TraceableError = (*traceableErrorImpl)(nil)

type traceableErrorImpl struct {
	Code    string
	Message string
	Err     error
	Origin  string
}

func (e traceableErrorImpl) GetMessage() string { return e.Message }

func (e traceableErrorImpl) WithError(err error) TraceableError {
	e.Err = err
	return e
}

func (e traceableErrorImpl) WithMessage(message string) TraceableError {
	e.Message = message
	return e
}

func (e traceableErrorImpl) WithOrigin() TraceableError {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return e
	}

	e.Origin = fmt.Sprintf("%s:%d", file, line)
	return e
}

func (e traceableErrorImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id      string `json:"id"`
		Message string `json:"message"`
	}{
		Id:      e.Code,
		Message: e.Message,
	})
}

func (e traceableErrorImpl) Error() string {
	asString := e.Code
	if e.Message != "" {
		asString += " - " + e.Message
	}
	if e.Err != nil {
		asString += ": " + e.Err.Error()
	}
	if e.Origin != "" {
		asString += " at " + e.Origin
	}
	return asString
}

func (e traceableErrorImpl) Unwrap() error { return e.Err }

func (e traceableErrorImpl) Is(target error) bool {
	targetErr, ok := target.(traceableErrorImpl)
	if !ok {
		return false
	}

	return e.Code == targetErr.Code
}
