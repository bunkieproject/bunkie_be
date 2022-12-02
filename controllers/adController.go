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
	}
}
