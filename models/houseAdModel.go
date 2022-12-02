package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type HouseAd struct {
	ID              primitive.ObjectID `bson:"_id"`
	User_id         *string            `bson:"user_id" validate:"required"`
	Header          *string            `json:"header" validate:"required,min=2,max=50"`
	Description     *string            `json:"description" validate:"required,min=2,max=100"`
	School          *string            `json:"school" validate:"required"`
	City            *string            `json:"city" validate:"required"`
	District        *string            `json:"district" validate:"required"`
	Neighborhood    *string            `json:"neighborhood" validate:"required"`
	Price           *string            `json:"price" validate:"required"`
	Size            *string            `json:"size" validate:"required"`
	Number_of_rooms *string            `json:"number_of_rooms" validate:"required"`
	House_photo     *string            `json:"house_photo" validate:"required"`
}

type UpdatedHouseAd struct {
	ID              primitive.ObjectID `bson:"_id"`
	User_id         *string            `bson:"user_id" validate:"required"`
	Header          *string            `json:"header" validate:"min=2,max=50`
	Description     *string            `json:"description" validate:"min=2,max=100"`
	School          *string            `json:"school"`
	City            *string            `json:"city"`
	District        *string            `json:"district"`
	Neighborhood    *string            `json:"neighborhood"`
	Price           *string            `json:"price"`
	Size            *string            `json:"size"`
	Number_of_rooms *string            `json:"number_of_rooms"`
	House_photo     *string            `json:"house_photo"`
}
