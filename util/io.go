package util

import (
	"io"
)

func ReadAllThenClose[T io.ReadCloser](r T) ([]byte, error) {
	defer r.Close()
	return io.ReadAll(r)
}
