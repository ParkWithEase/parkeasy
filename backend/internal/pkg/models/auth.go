package models

var (
	ErrAuthEmailExists     = CodeDuplicate.WithMsg("user with given email already exists")
	ErrAuthEmailOrPassword = CodeInvalidCredentials.WithMsg("invalid email or password")
	ErrRegInvalidEmail     = CodeInvalidEmail.WithMsg("email is invalid")
	ErrRegPasswordLength   = CodePasswordLength.WithMsg("password is too long or too short")
	ErrResetTokenInvalid   = CodeInvalidCredentials.WithMsg("password reset token invalid")
)

type EmailPasswordLoginInput struct {
	Email    string `json:"email" format:"email" doc:"User's email address"`
	Password string `json:"password" doc:"User's password"`
	Persist  bool   `json:"persist,omitempty" default:"false" doc:"Whether the resulting session should be persistent"`
}

// Hashed password ready for storage
type HashedPassword []byte
