package utils

import (
	"context"
	"errors"
)

type err struct {
	msg  string
	code uint
}

func (e err) Error() string {
	return e.msg
}

const (
	InternalErr = iota
	BadRequestErr
	NotFoundErr
)

func NewError(msg string, code uint) error {
	return &err{
		msg:  msg,
		code: code,
	}
}

const (
	internalErrorMsg = "internal error"
)

func FromError(ctx context.Context, in error) (string, uint) {
	var e *err
	if !errors.As(in, &e) {
		return internalErrorMsg, InternalErr
	}

	return e.msg, e.code
}
