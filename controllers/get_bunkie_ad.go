package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/bunkieproject/bunkie_be/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

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
