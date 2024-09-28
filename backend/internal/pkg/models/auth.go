package models

var (
	ErrAuthEmailExists     = NewUserFacingError("user with given email already exists")
	ErrAuthEmailOrPassword = NewUserFacingError("invalid email or password")
	ErrRegInvalidEmail     = NewUserFacingError("email is invalid")
	ErrRegPasswordLength   = NewUserFacingError("password is too long or too short")
	ErrResetTokenInvalid   = NewUserFacingError("password reset token invalid")
)

type EmailPasswordLoginInput struct {
	Email    string `json:"email" format:"email" doc:"User's email address"`
	Password string `json:"password" doc:"User's password"`
	Persist  bool   `json:"persist,omitempty" default:"false" doc:"Whether the resulting session should be persistent"`
}

// Hashed password ready for storage
type HashedPassword []byte
