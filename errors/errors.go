// Package errors provides a way to return detailed information
// for an RPC request error. The error is normally JSON encoded.
package errors // import "go.unistack.org/micro/v3/errors"

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var (
	// ErrBadRequest returns then requests contains invalid data
	ErrBadRequest = &Error{Code: 400}
	// ErrUnauthorized returns then user have unauthorized call
	ErrUnauthorized = &Error{Code: 401}
	// ErrForbidden returns then user have not access the resource
	ErrForbidden = &Error{Code: 403}
	// ErrNotFound returns then user specify invalid endpoint
	ErrNotFound = &Error{Code: 404}
	// ErrMethodNotAllowed returns then user try to get invalid method
	ErrMethodNotAllowed = &Error{Code: 405}
	// ErrTimeout returns then timeout exceeded
	ErrTimeout = &Error{Code: 408}
	// ErrConflict returns then request create duplicate resource
	ErrConflict = &Error{Code: 409}
	// ErrInternalServerError returns then server cant process request because of internal error
	ErrInternalServerError = &Error{Code: 500}
	// ErNotImplemented returns then server does not have desired endpoint method
	ErNotImplemented = &Error{Code: 501}
	// ErrBadGateway returns then server cant process request
	ErrBadGateway = &Error{Code: 502}
	// ErrServiceUnavailable returns then service unavailable
	ErrServiceUnavailable = &Error{Code: 503}
	// ErrGatewayTimeout returns then server have long time to process request
	ErrGatewayTimeout = &Error{Code: 504}
)

// Error type
type Error struct {
	// ID holds error id or service, usually someting like my_service or id
	ID string
	// Detail holds some useful details about error
	Detail string
	// Status usually holds text of http status
	Status string
	// Code holds error code
	Code int32
}

// Error satisfies error interface
func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

/*
// Generator struct holds id of error
type Generator struct {
	id string
}

// Generator can emit new error with static id
func NewGenerator(id string) *Generator {
	return &Generator{id: id}
}

func (g *Generator) BadRequest(format string, args ...interface{}) error {
	return BadRequest(g.id, format, args...)
}
*/

// New generates a custom error
func New(id, detail string, code int32) error {
	return &Error{
		ID:     id,
		Code:   code,
		Detail: detail,
		Status: http.StatusText(int(code)),
	}
}

// Parse tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func Parse(err string) *Error {
	e := &Error{}
	nerr := json.Unmarshal([]byte(err), e)
	if nerr != nil {
		e.Detail = err
	}
	return e
}

// BadRequest generates a 400 error.
func BadRequest(id, format string, args ...interface{}) error {
	return &Error{
		ID:     id,
		Code:   400,
		Detail: fmt.Sprintf(format, args...),
		Status: http.StatusText(400),
	}
}

// Unauthorized generates a 401 error.
func Unauthorized(id, format string, args ...interface{}) error {
	return &Error{
		ID:     id,
		Code:   401,
		Detail: fmt.Sprintf(format, args...),
		Status: http.StatusText(401),
	}
}

// Forbidden generates a 403 error.
func Forbidden(id, format string, args ...interface{}) error {
	return &Error{
		ID:     id,
		Code:   403,
		Detail: fmt.Sprintf(format, args...),
		Status: http.StatusText(403),
	}
}

// NotFound generates a 404 error.
func NotFound(id, format string, args ...interface{}) error {
	return &Error{
		ID:     id,
		Code:   404,
		Detail: fmt.Sprintf(format, args...),
		Status: http.StatusText(404),
	}
}

// MethodNotAllowed generates a 405 error.
func MethodNotAllowed(id, format string, args ...interface{}) error {
	return &Error{
		ID:     id,
		Code:   405,
		Detail: fmt.Sprintf(format, args...),
		Status: http.StatusText(405),
	}
}

// Timeout generates a 408 error.
func Timeout(id, format string, args ...interface{}) error {
	return &Error{
		ID:     id,
		Code:   408,
		Detail: fmt.Sprintf(format, args...),
		Status: http.StatusText(408),
	}
}

// Conflict generates a 409 error.
func Conflict(id, format string, args ...interface{}) error {
	return &Error{
		ID:     id,
		Code:   409,
		Detail: fmt.Sprintf(format, args...),
		Status: http.StatusText(409),
	}
}

// InternalServerError generates a 500 error.
func InternalServerError(id, format string, args ...interface{}) error {
	return &Error{
		ID:     id,
		Code:   500,
		Detail: fmt.Sprintf(format, args...),
		Status: http.StatusText(500),
	}
}

// NotImplemented generates a 501 error
func NotImplemented(id, format string, args ...interface{}) error {
	return &Error{
		ID:     id,
		Code:   501,
		Detail: fmt.Sprintf(format, args...),
		Status: http.StatusText(501),
	}
}

// BadGateway generates a 502 error
func BadGateway(id, format string, args ...interface{}) error {
	return &Error{
		ID:     id,
		Code:   502,
		Detail: fmt.Sprintf(format, args...),
		Status: http.StatusText(502),
	}
}

// ServiceUnavailable generates a 503 error
func ServiceUnavailable(id, format string, args ...interface{}) error {
	return &Error{
		ID:     id,
		Code:   503,
		Detail: fmt.Sprintf(format, args...),
		Status: http.StatusText(503),
	}
}

// GatewayTimeout generates a 504 error
func GatewayTimeout(id, format string, args ...interface{}) error {
	return &Error{
		ID:     id,
		Code:   504,
		Detail: fmt.Sprintf(format, args...),
		Status: http.StatusText(504),
	}
}

// Equal tries to compare errors
func Equal(err1 error, err2 error) bool {
	verr1, ok1 := err1.(*Error)
	verr2, ok2 := err2.(*Error)

	if ok1 != ok2 {
		return false
	}

	if !ok1 {
		return err1 == err2
	}

	if verr1.Code != verr2.Code {
		return false
	}

	return true
}

// FromError try to convert go error to *Error
func FromError(err error) *Error {
	if verr, ok := err.(*Error); ok && verr != nil {
		return verr
	}

	return Parse(err.Error())
}

// MarshalJSON returns error data
func (e *Error) MarshalJSON() ([]byte, error) {
	return e.Marshal()
}

// UnmarshalJSON set error data
func (e *Error) UnmarshalJSON(data []byte) error {
	return e.Unmarshal(data)
}

// ProtoMessage noop func
func (e *Error) ProtoMessage() {}

// Reset resets error
func (e *Error) Reset() {
	*e = Error{}
}

// String returns error as string
func (e *Error) String() string {
	return fmt.Sprintf(`{"id":"%s","detail":"%s","status":"%s","code":%d}`, e.ID, e.Detail, e.Status, e.Code)
}

// Marshal returns error data
func (e *Error) Marshal() ([]byte, error) {
	return []byte(e.String()), nil
}

// Unmarshal set error data
func (e *Error) Unmarshal(data []byte) error {
	str := string(data)
	if len(data) < 41 {
		return fmt.Errorf("invalid data")
	}
	parts := strings.FieldsFunc(str[1:len(str)-1], func(r rune) bool {
		return r == ','
	})
	for _, part := range parts {
		nparts := strings.FieldsFunc(part, func(r rune) bool {
			return r == ':'
		})
		for idx := 0; idx < len(nparts)/2; idx += 2 {
			val := strings.Trim(nparts[idx+1], `"`)
			if len(val) == 0 {
				continue
			}
			switch {
			case nparts[idx] == `"id"`:
				e.ID = val
			case nparts[idx] == `"detail"`:
				e.Detail = val
			case nparts[idx] == `"status"`:
				e.Status = val
			case nparts[idx] == `"code"`:
				c, err := strconv.ParseInt(val, 10, 32)
				if err != nil {
					return err
				}
				e.Code = int32(c)
			}
			idx++
		}
	}
	return nil
}
