package controllers

import (
	"net/http"

	"github.com/bunkieproject/bunkie_be/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// Logout logs out a user
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.LogoutRequest
		var onlineToken models.Token

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := onlineCollection.FindOne(c, bson.M{"token": request.Token}).Decode(&onlineToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found 1"})
			return
		}

		if onlineToken.User_id != request.User_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found 2"})
			return
		}

		resultDeleteNumber, err := onlineCollection.DeleteOne(c, bson.M{"token": request.Token})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": resultDeleteNumber})
	}
}
