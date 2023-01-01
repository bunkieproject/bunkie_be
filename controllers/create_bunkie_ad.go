package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/bunkieproject/bunkie_be/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
