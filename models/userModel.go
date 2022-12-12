package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountInfo struct {
	ID           primitive.ObjectID `bson:"_id"`
	Email        *string            `json:"email"`
	Username     *string            `json:"username"`
	Password     *string            `json:"Password"`
	PasswordHash *string            `json:"password_hash"`
	Token        *string            `json:"token"`
	UserType     *string            `json:"user_type"`
	User_id      string             `json:"user_id"`
	Is_banned    bool               `json:"is_banned" default:"false"`
	ProfileInfo  *ProfileInfo       `json:"profile_info"`
	CreatedAt    time.Time          `json:"created_at"`
}

type ProfileInfo struct {
	FullName       *string `json:"name"`
	Phone          *string `json:"phone"`
	ProfilePicture *string `json:"profile_picture"`
	Description    *string `json:"description"`
	Gender         *string `json:"gender"`
}
