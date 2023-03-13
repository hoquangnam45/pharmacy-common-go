package response

import (
	"net/http"
	"time"
)

var ErrInternalServer = NewResponseError(500, "internal server error")

type Error struct {
	Msg        string    `json:"msg"`
	Timestamp  time.Time `json:"timestamp"`
	Path       string    `json:"path"`
	StatusCode int       `json:"-"`
}

func NewResponseError(statusCode int, msg string) *Error {
	return &Error{
		Msg:        msg,
		StatusCode: statusCode,
	}
}

func (e *Error) Error() string {
	return e.Msg
}

func (e *Error) AddContext(r *http.Request) *Error {
	return &Error{
		Msg:        e.Msg,
		Timestamp:  time.Now(),
		Path:       r.URL.Path,
		StatusCode: e.StatusCode,
	}
}
