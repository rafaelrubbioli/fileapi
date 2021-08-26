package gqlerror

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var (
	ErrServiceUnavailable = newTyped("service unavailable", ServiceUnavailableType)
	ErrFileTooBig         = newTyped("max file size is 500b", BadRequestType)
	ErrInvalidPath        = newTyped("path cannot contain '..'", BadRequestType)
	ErrInvalidID          = newTyped("invalid id", BadRequestType)
	ErrNotYetSupported    = newTyped("not yet supported", ServiceUnavailableType)
	ErrNotFound           = newTyped("not found", NotFoundType)
)

type ErrorType string

const ErrCodeLabel = "code"

var (
	DuplicatedType         ErrorType = "DUPLICATED"
	NotFoundType           ErrorType = "NOT_FOUND"
	ServiceUnavailableType ErrorType = "SERVICE_UNAVAILABLE"
	UnauthorizedType       ErrorType = "UNAUTHORIZED"
	BadRequestType         ErrorType = "BAD_REQUEST"
)

var errorMap = map[error]error{}

func Error(err error) error {
	if newErr, ok := errorMap[err]; ok {
		return newErr
	}

	var errNoSuchKey *types.NoSuchKey
	if errors.As(err, &errNoSuchKey) {
		return ErrNotFound
	}

	var errNotFound *types.NotFound
	if errors.As(err, &errNotFound) {
		return errNotFound
	}

	log.Println(err)
	return ErrServiceUnavailable
}

func new(message string, params ...interface{}) *gqlerror.Error {
	err := &gqlerror.Error{
		Message: message,
	}

	if len(params) > 0 {
		err.Message = fmt.Sprintf(message, params...)
	}

	return err
}

func newTyped(message string, errorType ErrorType, params ...interface{}) *gqlerror.Error {
	err := new(message, params...)
	err.Extensions = map[string]interface{}{
		ErrCodeLabel: errorType,
	}

	return err
}
