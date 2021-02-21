package server

import "github.com/unistack-org/micro/v3/errors"

type Error struct {
	id string
}

func NewError(id string) *Error {
	return &Error{id}
}

func (e *Error) BadRequest(format string, a ...interface{}) error {
	return errors.BadRequest(e.id, format, a...)
}

func (e *Error) Unauthorized(format string, a ...interface{}) error {
	return errors.Unauthorized(e.id, format, a...)
}

func (e *Error) Forbidden(format string, a ...interface{}) error {
	return errors.Forbidden(e.id, format, a...)
}

func (e *Error) NotFound(format string, a ...interface{}) error {
	return errors.NotFound(e.id, format, a...)
}

func (e *Error) MethodNotAllowed(format string, a ...interface{}) error {
	return errors.MethodNotAllowed(e.id, format, a...)
}

func (e *Error) Timeout(format string, a ...interface{}) error {
	return errors.Timeout(e.id, format, a...)
}

func (e *Error) Conflict(format string, a ...interface{}) error {
	return errors.Conflict(e.id, format, a...)
}

func (e *Error) InternalServerError(format string, a ...interface{}) error {
	return errors.InternalServerError(e.id, format, a...)
}

func (e *Error) NotImplemented(format string, a ...interface{}) error {
	return errors.NotImplemented(e.id, format, a...)
}

func (e *Error) BadGateway(format string, a ...interface{}) error {
	return errors.BadGateway(e.id, format, a...)
}

func (e *Error) ServiceUnavailable(format string, a ...interface{}) error {
	return errors.ServiceUnavailable(e.id, format, a...)
}

func (e *Error) GatewayTimeout(format string, a ...interface{}) error {
	return errors.GatewayTimeout(e.id, format, a...)
}
