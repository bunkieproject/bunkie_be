package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BunkieAd struct {
	ID           primitive.ObjectID `bson:"_id"`
	User_id      *string            `bson:"user_id" validate:"required"`
	City         *string            `json:"city" validate:"required"`
	District     *string            `json:"district" validate:"required"`
	Neighborhood *string            `json:"neighborhood" validate:"required"`
	Size         *string            `json:"size" validate:"required"`
	Price        *string            `json:"price" validate:"required"`
	Gender       *string            `json:"gender" validate:"required"`
	Job          *string            `json:"job" validate:"required"`
}

type UpdatedBunkieAd struct {
	ID           primitive.ObjectID `bson:"_id"`
	User_id      *string            `bson:"user_id" validate:"required"`
	City         *string            `json:"city"`
	District     *string            `json:"district"`
	Neighborhood *string            `json:"neighborhood"`
	Size         *string            `json:"size"`
	Price        *string            `json:"price"`
	Gender       *string            `json:"gender"`
	Job          *string            `json:"job"`
}
