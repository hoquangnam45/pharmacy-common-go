package util

type GroupError struct {
	Cause error
	Group error
}

func NewGroupError(cause error, group error) *GroupError {
	return &GroupError{
		Cause: cause,
		Group: group,
	}
}

func (e *GroupError) Unwrap() error {
	return e.Group
}

func (e *GroupError) Error() string {
	return e.Group.Error()
}
