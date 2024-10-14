package models

var ErrNoProfile = CodeNoProfile.WithMsg("no profile exists for this user")

type UserProfile struct {
	// The full name of an user
	FullName string `json:"full_name" doc:"The user's full name"`
	// The email of that user
	Email string `json:"email" format:"email" doc:"The user email address"`
}

type UserCreationInput struct {
	UserProfile
	Password string `json:"password" doc:"The user password"`
}
