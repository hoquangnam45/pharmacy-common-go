package request

type Error struct {
	*Response
}

func ToErrorResponse(r *Response) *Error {
	return &Error{r}
}

func (e *Error) Error() string {
	return e.Msg
}

func NewErrorResponse(msg string, statusCode int) *Error {
	return &Error{
		Response: &Response{
			Msg:        msg,
			StatusCode: statusCode,
		},
	}
}
