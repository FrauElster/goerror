package goerror

import (
	"encoding/json"
	"errors"
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

// WithError sets the err to the TraceableError
func (e traceableErrorImpl) WithError(err error) TraceableError {
	e.Err = errors.Join(e.Err, err)
	return e
}

// WithMessage sets the message to the TraceableError
func (e traceableErrorImpl) WithMessage(message string) TraceableError {
	e.Message = message
	return e
}

// WithOrigin applies the calles file location to TraceableError
func (e traceableErrorImpl) WithOrigin() TraceableError {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return e
	}

	e.Origin = fmt.Sprintf("%s:%d", file, line)
	return e
}

// MarshalJSON implements the Marshaller interface
func (e traceableErrorImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id      string `json:"id"`
		Message string `json:"message"`
	}{
		Id:      e.Code,
		Message: e.Message,
	})
}

// Error implements the Error interface.
// It returns a composition of Error.Code, message (if available), underlying error (if available), and origin (if available)
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

// Is implements errors.Is interface, to check equality on error Id
func (e traceableErrorImpl) Is(target error) bool {
	targetErr, ok := target.(traceableErrorImpl)
	if !ok {
		return false
	}

	return e.Code == targetErr.Code
}
