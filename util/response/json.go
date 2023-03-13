package response

import (
	"encoding/json"
	"net/http"
)

type Json[T any] struct {
	resp       T
	statusCode int
}

func NewJsonResponse[T any](statusCode int, resp T) *Json[T] {
	return &Json[T]{resp, statusCode}
}

func NewErrorJsonResponse(r *http.Request, err *Error) *Json[*Error] {
	return NewJsonResponse(err.StatusCode, err.AddContext(r))
}

func (j *Json[T]) Response() (string, error) {
	if data, err := json.Marshal(j.resp); err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

func (j *Json[T]) StatusCode() int {
	return j.statusCode
}

func (j *Json[T]) ContentType() string {
	return "application/json"
}
