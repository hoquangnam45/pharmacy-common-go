package request

import "github.com/hoquangnam45/pharmacy-common-go/util"

type Json[T any] struct {
	*Response
}

func ToJsonResponse[T any](r *Response) *Json[T] {
	return &Json[T]{r}
}

func (j *Json[T]) Get(placeholder *T) (*T, error) {
	return util.UnmarshalJson(placeholder)([]byte(j.Msg))
}
