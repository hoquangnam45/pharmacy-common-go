package request

import (
	"io"
	"net/http"
)

type Response struct {
	Msg        string
	StatusCode int
}

func NewRequestResponse(r *http.Response) (*Response, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return &Response{
		Msg:        string(data),
		StatusCode: r.StatusCode,
	}, nil
}
