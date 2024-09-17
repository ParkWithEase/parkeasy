package models

type Greeting struct {
	Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
}
