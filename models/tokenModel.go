package models

type Token struct {
	Token   *string `json:"token" validate:"required"`
	User_id *string `json:"user_id" validate:"required"`
}
