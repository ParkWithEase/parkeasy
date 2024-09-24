package models

// Errors meant to be sent to clients
//
// This allows the router to automatically filter server errors and prevent them
// from reaching clients.
type UserFacingError struct {
	msg string
}

func NewUserFacingError(msg string) error {
	return &UserFacingError{msg}
}

// Implements the `error` protocol
func (e *UserFacingError) Error() string {
	return e.msg
}
