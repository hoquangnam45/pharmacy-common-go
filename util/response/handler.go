package response

import (
	"errors"
	"net/http"
)

func Handler[T any](response *Json[T], w http.ResponseWriter) error {
	resp, err := response.Response()
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", response.ContentType())
	w.WriteHeader(response.StatusCode())
	w.Write([]byte(resp))
	return nil
}

func ErrorHandler(err error, r *http.Request, w http.ResponseWriter) error {
	apiErr := &Error{}
	if errors.As(err, &apiErr) {
		jsonResponse := NewErrorJsonResponse(r, apiErr)
		resp, err := jsonResponse.Response()
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", jsonResponse.ContentType())
		w.WriteHeader(jsonResponse.StatusCode())
		w.Write([]byte(resp))
		return nil
	} else {
		jsonResponse := NewErrorJsonResponse(r, ErrInternalServer)
		resp, err := jsonResponse.Response()
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", jsonResponse.ContentType())
		w.WriteHeader(jsonResponse.StatusCode())
		w.Write([]byte(resp))
		return nil
	}
}
