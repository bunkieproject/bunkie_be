package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/bunkieproject/bunkie_be/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

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
