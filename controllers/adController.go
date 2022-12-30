package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/bunkieproject/bunkie_be/database"
	"github.com/bunkieproject/bunkie_be/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var bunkieAdCollection *mongo.Collection = database.OpenCollection(database.Client, "bunkie_ads")
var roomAdCollection *mongo.Collection = database.OpenCollection(database.Client, "room_ads")

func CreateBunkieAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
		var ad models.BunkieAd
		var request models.CreateBunkieRequest
		var user models.AccountInfo

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !checkIfUserOnline(request.User_id, c) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not online"})
			return
		}

		validationErr := validate.Struct(ad)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"token": request.Token}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.User_id != request.User_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User id does not match token"})
			return
		}

		count_ads, err := bunkieAdCollection.CountDocuments(ctx, bson.M{"user_id": request.User_id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count_ads > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User can only have one bunkie ad"})
			return
		}

		ad = models.BunkieAd{
			Ad_id:           primitive.NewObjectID().Hex(),
			User_id:         request.User_id,
			School:          request.School,
			City:            request.City,
			District:        request.District,
			Quarter:         request.Quarter,
			Header:          request.Header,
			Description:     request.Description,
			NumberOfRooms:   request.NumberOfRooms,
			Price:           request.Price,
			GenderPreferred: request.GenderPreferred,
			CreatedAt:       time.Now(),
		}

		_, err = bunkieAdCollection.InsertOne(ctx, ad)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ad)
	}
}

func GetBunkieAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
		var ad models.BunkieAd
		var request models.GetBunkieAdRequest

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !checkIfUserOnline(request.User_id, c) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not online"})
			return
		}

		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := bunkieAdCollection.FindOne(ctx, bson.M{"ad_id": *request.Ad_id}).Decode(&ad)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ad)
	}
}

func UpdateBunkieAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
		var ad models.BunkieAd
		var request models.UpdateBunkieRequest
		var user models.AccountInfo

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !checkIfUserOnline(request.User_id, c) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not online"})
			return
		}

		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"token": request.Token}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.User_id != request.User_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User id does not match token"})
			return
		}

		err = bunkieAdCollection.FindOne(ctx, bson.M{"ad_id": request.Ad_id}).Decode(&ad)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		update := getValidUpdatesB(request)

		_, err = bunkieAdCollection.UpdateOne(ctx, bson.M{"ad_id": request.Ad_id}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Bunkie ad updated successfully"})
	}
}

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

func DeleteBunkieAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
		var request models.DeleteBunkieRequest
		var user models.AccountInfo

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !checkIfUserOnline(request.User_id, c) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not online"})
			return
		}

		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"token": request.Token}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.User_id != request.User_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User id does not match token"})
			return
		}

		_, err = bunkieAdCollection.DeleteOne(ctx, bson.M{"ad_id": request.Ad_id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Bunkie ad deleted successfully"})
	}
}

func CreateRoomAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
		var ad models.RoomAd
		var user models.AccountInfo
		var request models.CreateRoomAdRequest

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !checkIfUserOnline(request.User_id, c) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not online"})
			return
		}

		validationErr := validate.Struct(ad)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"token": request.Token}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.User_id != request.User_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User id and token do not match"})
			return
		}

		ad = models.RoomAd{
			Ad_id:            primitive.NewObjectID().Hex(),
			User_id:          request.User_id,
			Header_bytearray: request.Header_bytearray,
			Other_bytearrays: request.Other_bytearrays,
			School:           request.School,
			City:             request.City,
			District:         request.District,
			Quarter:          request.Quarter,
			Header:           request.Header,
			Description:      request.Description,
			NumberOfRooms:    request.NumberOfRooms,
			Price:            request.Price,
			GenderPreferred:  request.GenderPreferred,
			CreatedAt:        time.Now(),
		}

		_, err = roomAdCollection.InsertOne(ctx, ad)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ad)
	}
}

func GetRoomAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
		var ad models.RoomAd
		var request models.GetRoomAdRequest

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !checkIfUserOnline(request.User_id, c) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not online"})
			return
		}

		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := roomAdCollection.FindOne(ctx, bson.M{"ad_id": request.Ad_id}).Decode(&ad)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ad)
	}
}

func GetRoomAds() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
		var ads []models.RoomAd
		var request models.GetRoomAdsRequest

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !checkIfUserOnline(request.User_id, c) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not online"})
			return
		}

		err := validate.Struct(request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cursor, err := roomAdCollection.Find(ctx, bson.M{"user_id": request.User_id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for cursor.Next(ctx) {
			var ad models.RoomAd
			cursor.Decode(&ad)
			ads = append(ads, ad)
		}

		c.JSON(http.StatusOK, ads)
	}
}

func UpdateRoomAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
		var ad models.RoomAd
		var user models.AccountInfo
		var request models.UpdateRoomAdRequest

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !checkIfUserOnline(request.User_id, c) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not online"})
			return
		}

		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"token": request.Token}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.User_id != request.User_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User id does not match token"})
			return
		}

		err = roomAdCollection.FindOne(ctx, bson.M{"ad_id": request.Ad_id}).Decode(&ad)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if ad.User_id != request.User_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User id does not match ad owner"})
			return
		}

		update := getValidUpdatesR(request)

		_, err = roomAdCollection.UpdateOne(ctx, bson.M{"ad_id": request.Ad_id}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Room ad updated"})
	}
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

func DeleteRoomAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
		var ad models.RoomAd
		var user models.AccountInfo
		var request models.DeleteRoomAdRequest

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !checkIfUserOnline(request.User_id, c) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not online"})
			return
		}

		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"token": request.Token}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.User_id != request.User_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User id does not match token"})
			return
		}

		err = roomAdCollection.FindOne(ctx, bson.M{"ad_id": request.Ad_id}).Decode(&ad)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if ad.User_id != request.User_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User id does not match ad owner"})
			return
		}

		_, err = roomAdCollection.DeleteOne(ctx, bson.M{"ad_id": request.Ad_id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Room ad deleted"})
	}
}

func SearchBunkieAdDefault() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.AccountInfo
		var request models.SearchBunkieDefaultRequest

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !checkIfUserOnline(request.User_id, c) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not online"})
			return
		}

		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"token": request.Token}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.User_id != request.User_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User id does not match token"})
			return
		}

		opts := options.Find().SetLimit(request.HowManyDocs)
		opts.SetSort(bson.D{{"created_at", -1}})
		var filter bson.M
		if user.ProfileInfo != nil {
			if user.ProfileInfo.City != nil {
				filter = bson.M{"city": *user.ProfileInfo.City}
			} else {
				filter = bson.M{"city": "İstanbul"}
			}
		} else {
			filter = bson.M{"city": "İstanbul"}
		}

		cursor, err := bunkieAdCollection.Find(ctx, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var results []models.BunkieAd
		for cursor.Next(ctx) {
			var ad models.BunkieAd
			cursor.Decode(&ad)
			results = append(results, ad)
		}

		c.JSON(http.StatusOK, results)
	}
}

func SearchBunkieAdPreferred() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.AccountInfo
		var request models.SearchBunkiePreferredRequest

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !checkIfUserOnline(request.User_id, c) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not online"})
			return
		}

		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"token": request.Token}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.User_id != request.User_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User id does not match token"})
			return
		}

		opts := options.Find().SetLimit(request.HowManyDocs)
		opts.SetSort(bson.D{{"created_at", -1}})

		cursor, err := bunkieAdCollection.Find(ctx, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var results []models.BunkieAd
		for cursor.Next(ctx) {
			var ad models.BunkieAd
			cursor.Decode(&ad)
			results = append(results, ad)
		}

		c.JSON(http.StatusOK, results)

	}
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

func SearchRoomAdDefault() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.AccountInfo
		var request models.SearchRoomAdDefaultRequest

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !checkIfUserOnline(request.User_id, c) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not online"})
			return
		}

		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"token": request.Token}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.User_id != request.User_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User id does not match token"})
			return
		}

		opts := options.Find().SetLimit(request.HowManyDocs)
		opts.SetSort(bson.D{{"created_at", -1}})
		var filter bson.M
		if user.ProfileInfo != nil {
			if user.ProfileInfo.City != nil {
				filter = bson.M{"city": *user.ProfileInfo.City}
			} else {
				filter = bson.M{"city": "İstanbul"}
			}
		} else {
			filter = bson.M{"city": "İstanbul"}
		}

		cursor, err := roomAdCollection.Find(ctx, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var results []models.RoomAd
		for cursor.Next(ctx) {
			var ad models.RoomAd
			cursor.Decode(&ad)
			results = append(results, ad)
		}

		c.JSON(http.StatusOK, results)
	}
}

func SearchRoomAdPreferred() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.AccountInfo
		var request models.SearchRoomAdPreferredRequest

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !checkIfUserOnline(request.User_id, c) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not online"})
			return
		}

		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"token": request.Token}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.User_id != request.User_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User id does not match token"})
			return
		}

		opts := options.Find().SetLimit(request.HowManyDocs)
		opts.SetSort(bson.D{{"created_at", -1}})
		filter := getAppropriateFilterRoom(request)

		cursor, err := roomAdCollection.Find(ctx, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var results []models.RoomAd
		for cursor.Next(ctx) {
			var ad models.RoomAd
			cursor.Decode(&ad)
			results = append(results, ad)
		}

		c.JSON(http.StatusOK, results)

	}
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
