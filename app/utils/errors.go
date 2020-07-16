package utils

import (
	"net/http"
)

type APIError interface {
	error
	Status() int
}

type BaseAPIError struct {
	Msg string
	StatusCode int
}

func (be BaseAPIError) Error() string {
	return be.Msg
}

func (be BaseAPIError) Status() int {
	return be.StatusCode
}


func NewInternalError(msg string) *BaseAPIError {
	return &BaseAPIError{
		Msg: msg,
		StatusCode: 500,
	}
}

func NewIllegalArgument(msg string) *BaseAPIError {
	return &BaseAPIError{
		Msg: msg,
		StatusCode: 400,
	}
}

func WriteHttpError(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(APIError); ok {
		http.Error(w, apiErr.Error(), apiErr.Status())
	} else {
		http.Error(w, err.Error(), 500)
	}
}