package models

import (
	"time"
)

type BunkieAd struct {
	Ad_id           string    `json:"ad_id"`
	User_id         string    `json:"user_id"`
	City            *string   `json:"city"`
	District        *string   `json:"district"`
	Quarter         *string   `json:"quarter"`
	Header          *string   `json:"header"`
	Description     *string   `json:"description"`
	NumberOfRooms   *string   `json:"number_of_rooms"`
	Price           float64   `json:"price"`
	GenderPreferred *string   `json:"gender_preferred"`
	CreatedAt       time.Time `json:"created_at"`
}

type RoomAd struct {
	Ad_id            string    `json:"ad_id"`
	User_id          string    `json:"user_id"`
	Header_bytearray *string   `json:"header_bytearray"`
	Other_bytearrays *string   `json:"other_bytearrays"`
	Header           *string   `json:"header"`
	Description      *string   `json:"description"`
	City             *string   `json:"city"`
	District         *string   `json:"district"`
	Quarter          *string   `json:"quarter"`
	NumberOfRooms    *string   `json:"number_of_rooms"`
	Price            float64   `json:"price"`
	GenderPreferred  *string   `json:"gender_preferred"`
	CreatedAt        time.Time `json:"created_at"`
}
