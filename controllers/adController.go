package controllers

import (
	"bunkie_be/database"
	"bunkie_be/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var houseAdCollection *mongo.Collection = database.OpenCollection(database.Client, "houseAds")
var bunkieAdCollection *mongo.Collection = database.OpenCollection(database.Client, "bunkieAds")

func CreateHouseAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var houseAd models.HouseAd

		// Bind JSON to struct
		if err := c.BindJSON(&houseAd); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the input
		validationErr := validate.Struct(houseAd)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if user has already created an house ad
		count_ads, err := houseAdCollection.CountDocuments(ctx, bson.M{"user_id": houseAd.User_id})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count_ads > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already has an ad"})
			return
		}

		// Insert house ad into database
		_, err = houseAdCollection.InsertOne(ctx, houseAd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, houseAd)
	}
}

func DeleteHouseAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var houseAd models.UpdatedHouseAd
		var foundAd models.HouseAd

		// Bind JSON to struct
		if err := c.BindJSON(&houseAd); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := houseAdCollection.FindOne(ctx, bson.M{"user_id": houseAd.User_id}).Decode(&foundAd)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if houseAd.Header == nil {
			houseAd.Header = foundAd.Header
		}
		if houseAd.Description == nil {
			houseAd.Description = foundAd.Description
		}

		// Validate the input
		validationErr := validate.Struct(houseAd)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Delete house ad from database
		_, err = houseAdCollection.DeleteOne(ctx, bson.D{{Key: "user_id", Value: houseAd.User_id}})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

func UpdateHouseAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var houseAd models.UpdatedHouseAd
		var foundAd models.HouseAd

		// Bind JSON to struct
		if err := c.BindJSON(&houseAd); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := houseAdCollection.FindOne(ctx, bson.M{"user_id": houseAd.User_id}).Decode(&foundAd)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if houseAd.Header == nil {
			houseAd.Header = foundAd.Header
		}
		if houseAd.Description == nil {
			houseAd.Description = foundAd.Description
		}

		// Validate the input
		validationErr := validate.Struct(houseAd)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		update_map := setUpdateMap(houseAd)

		update := bson.M{"$set": update_map}
		err = houseAdCollection.FindOneAndUpdate(ctx, bson.M{"user_id": houseAd.User_id}, update).Decode(&foundAd)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundAd)

	}
}

func setUpdateMap(ad models.UpdatedHouseAd) bson.M {
	update_map := bson.M{}
	if ad.Header != nil {
		update_map["header"] = ad.Header
	}
	if ad.Description != nil {
		update_map["description"] = ad.Description
	}
	if ad.School != nil {
		update_map["school"] = ad.School
	}
	if ad.City != nil {
		update_map["city"] = ad.City
	}
	if ad.District != nil {
		update_map["district"] = ad.District
	}
	if ad.Neighborhood != nil {
		update_map["neighborhood"] = ad.Neighborhood
	}
	if ad.Price != nil {
		update_map["price"] = ad.Price
	}
	if ad.Size != nil {
		update_map["size"] = ad.Size
	}
	if ad.Number_of_rooms != nil {
		update_map["number_of_rooms"] = ad.Number_of_rooms
	}
	if ad.House_photo != nil {
		update_map["house_photo"] = ad.House_photo
	}
	return update_map
}

func CreateBunkieAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var bunkieAd models.BunkieAd

		// Bind JSON to struct
		if err := c.BindJSON(&bunkieAd); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the input
		validationErr := validate.Struct(bunkieAd)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if user has already created a bunkie ad
		count_ads, err := bunkieAdCollection.CountDocuments(ctx, bson.M{"user_id": bunkieAd.User_id})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count_ads > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already has an ad"})
			return
		}

		// Create bunkie ad in database
		_, err = bunkieAdCollection.InsertOne(ctx, bunkieAd)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, bunkieAd)

	}
}

func DeleteBunkieAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var bunkieAd models.UpdatedBunkieAd

		// Bind JSON to struct
		if err := c.BindJSON(&bunkieAd); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the input
		validationErr := validate.Struct(bunkieAd)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		_, err := bunkieAdCollection.DeleteOne(ctx, bson.M{"user_id": bunkieAd.User_id})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Bunkie ad deleted"})
	}
}

func UpdateBunkieAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var bunkieAd models.UpdatedBunkieAd
		var foundAd models.BunkieAd

		// Bind JSON to struct
		if err := c.BindJSON(&bunkieAd); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the input
		validationErr := validate.Struct(bunkieAd)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if user has already created a bunkie ad
		count_ads, err := bunkieAdCollection.CountDocuments(ctx, bson.M{"user_id": bunkieAd.User_id})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count_ads == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not have an ad"})
			return
		}

		// Create update map
		update_map := setUpdateMapBunkie(bunkieAd)

		// Update bunkie ad in database
		err = bunkieAdCollection.FindOneAndUpdate(ctx, bson.M{"user_id": bunkieAd.User_id}, bson.M{"$set": update_map}).Decode(&foundAd)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundAd)
	}
}

func setUpdateMapBunkie(ad models.UpdatedBunkieAd) bson.M {
	update_map := bson.M{}
	if ad.City != nil {
		update_map["city"] = ad.City
	}
	if ad.District != nil {
		update_map["district"] = ad.District
	}
	if ad.Neighborhood != nil {
		update_map["neighborhood"] = ad.Neighborhood
	}
	if ad.Size != nil {
		update_map["size"] = ad.Size
	}
	if ad.Price != nil {
		update_map["price"] = ad.Price
	}
	if ad.Gender != nil {
		update_map["gender"] = ad.Gender
	}
	if ad.Job != nil {
		update_map["job"] = ad.Job
	}
	return update_map
}
