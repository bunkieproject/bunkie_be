package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountInfo struct {
	ID               primitive.ObjectID `bson:"_id"`
	First_name       *string            `json:"name" validate:"required,min=2,max=100"`
	Last_name        *string            `json:"surname" validate:"required,min=2,max=100"`
	Email            *string            `json:"email" validate:"email,required"`
	Phone            *string            `json:"phone" validate:"required"`
	Username         *string            `json:"username" validate:"required"`
	Password         *string            `json:"Password" validate:"required,min=6"`
	Password_confirm *string            `json:"Password_confirm" validate:"required,min=6"`
	Token            *string            `json:"token"`
	User_type        *string            `json:"user_type" validate:"required,eq=admin|eq=registeredUser"`
	User_id          string             `json:"user_id"`
	Is_banned        bool               `json:"is_banned" default:"false"`
	Is_online        bool               `json:"is_online" default:"false"`
}

type ProfileInfo struct {
	ID              primitive.ObjectID `bson:"_id"`
	Description     *string            `json:"description" validate:"required,min=2,max=200"`
	Profile_picture *string            `json:"profile_picture" validate:"required"`
	Gender          *int               `json:"gender"`
}
