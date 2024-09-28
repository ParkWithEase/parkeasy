package models

type PasswordUpdateInput struct {
	OldPassword string `json:"old_password" doc:"The user previous password"`
	NewPassword string `json:"new_password" doc:"The user new password"`
}

type PasswordResetInput struct {
	PasswordResetToken string `json:"password_token" doc:"The token used to reset user password"`
	NewPassword        string `json:"new_password" doc:"The user new password"`
}

type PasswordResetTokenRequest struct {
	Email string `json:"email" doc:"The user email"`
}
