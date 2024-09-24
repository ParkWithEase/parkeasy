package models

var (
	ErrAuthEmailExists     = NewUserFacingError("user with given email already exists")
	ErrAuthEmailOrPassword = NewUserFacingError("invalid email or password")
)

type EmailPasswordLoginInput struct {
	Email    string `json:"email" format:"email" doc:"User's email address"`
	Password string `json:"password" doc:"User's password"`
	Persist  bool   `json:"persist,omitempty" default:"false" doc:"Whether the resulting session should be persistent"`
}
