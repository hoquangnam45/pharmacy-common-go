package util

type GroupError struct {
	Cause error
	Group error
}

func NewGroupError(group error, cause error) *GroupError {
	return &GroupError{
		Cause: cause,
		Group: group,
	}
}

func (e *GroupError) Unwrap() error {
	return e.Cause
}

func (e *GroupError) Error() string {
	return e.Cause.Error()
}
