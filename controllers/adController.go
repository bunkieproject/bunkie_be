package controllers

import (
	"github.com/bunkieproject/bunkie_be/database"
	"github.com/bunkieproject/bunkie_be/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var bunkieAdCollection *mongo.Collection = database.OpenCollection(database.Client, "bunkie_ads")
var roomAdCollection *mongo.Collection = database.OpenCollection(database.Client, "room_ads")

func getValidUpdatesB(request models.UpdateBunkieRequest) bson.M {
	update := bson.M{"$set": bson.M{}}

	if request.School != nil {
		update["$set"].(bson.M)["school"] = request.School
	}
	if request.City != nil {
		update["$set"].(bson.M)["city"] = request.City
	}
	if request.District != nil {
		update["$set"].(bson.M)["district"] = request.District
	}
	if request.Quarter != nil {
		update["$set"].(bson.M)["quarter"] = request.Quarter
	}
	if request.Header != nil {
		update["$set"].(bson.M)["header"] = request.Header
	}
	if request.Description != nil {
		update["$set"].(bson.M)["description"] = request.Description
	}
	if request.NumberOfRooms != nil {
		update["$set"].(bson.M)["number_of_rooms"] = request.NumberOfRooms
	}
	if request.Price != 0 {
		update["$set"].(bson.M)["price"] = request.Price
	}
	if request.GenderPreferred != nil {
		update["$set"].(bson.M)["gender_preferred"] = request.GenderPreferred
	}
	return update
}

func getValidUpdatesR(request models.UpdateRoomAdRequest) bson.M {
	update := bson.M{"$set": bson.M{}}

	if request.Header_bytearray != nil {
		update["$set"].(bson.M)["header_bytearray"] = request.Header_bytearray
	}
	if request.Other_bytearrays != nil {
		update["$set"].(bson.M)["other_bytearrays"] = request.Other_bytearrays
	}
	if request.School != nil {
		update["$set"].(bson.M)["school"] = request.School
	}
	if request.City != nil {
		update["$set"].(bson.M)["city"] = request.City
	}
	if request.District != nil {
		update["$set"].(bson.M)["district"] = request.District
	}
	if request.Quarter != nil {
		update["$set"].(bson.M)["quarter"] = request.Quarter
	}
	if request.Header != nil {
		update["$set"].(bson.M)["header"] = request.Header
	}
	if request.Description != nil {
		update["$set"].(bson.M)["description"] = request.Description
	}
	if request.NumberOfRooms != nil {
		update["$set"].(bson.M)["number_of_rooms"] = request.NumberOfRooms
	}
	if request.Price != 0 {
		update["$set"].(bson.M)["price"] = request.Price
	}
	if request.GenderPreferred != nil {
		update["$set"].(bson.M)["gender_preferred"] = request.GenderPreferred
	}

	return update

}

func getAppropriateFilter(request models.SearchBunkiePreferredRequest) bson.M {
	filter := bson.M{}

	if request.LowerPrice != 0 {
		filter["price"] = bson.M{"$gte": request.LowerPrice}
	}
	if request.UpperPrice != 0 {
		filter["price"] = bson.M{"$lte": request.UpperPrice}
	}
	if request.GenderPreferred != nil {
		filter["gender_preferred"] = *request.GenderPreferred
	}
	if request.NumberOfRooms != nil {
		filter["number_of_rooms"] = *request.NumberOfRooms
	}
	if request.School != nil {
		filter["school"] = *request.School
	}
	if request.City != nil {
		filter["city"] = *request.City
	}
	if request.District != nil {
		filter["district"] = *request.District
	}
	if request.Quarter != nil {
		filter["quarter"] = *request.Quarter
	}

	return filter
}

func getAppropriateFilterRoom(request models.SearchRoomAdPreferredRequest) bson.M {
	filter := bson.M{}

	if request.LowerPrice != 0 {
		filter["price"] = bson.M{"$gte": request.LowerPrice}
	}
	if request.UpperPrice != 0 {
		filter["price"] = bson.M{"$lte": request.UpperPrice}
	}
	if request.GenderPreferred != nil {
		filter["gender_preferred"] = *request.GenderPreferred
	}
	if request.NumberOfRooms != nil {
		filter["number_of_rooms"] = *request.NumberOfRooms
	}
	if request.School != nil {
		filter["school"] = *request.School
	}
	if request.City != nil {
		filter["city"] = *request.City
	}
	if request.District != nil {
		filter["district"] = *request.District
	}
	if request.Quarter != nil {
		filter["quarter"] = *request.Quarter
	}

	return filter
}
