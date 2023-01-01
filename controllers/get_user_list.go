package controllers

import (
	"net/http"

	"github.com/bunkieproject/bunkie_be/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.GetUserListRequest
		var user models.AccountInfo
		var users []models.AccountInfo

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := userCollection.FindOne(c, bson.M{"token": request.Token}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if *user.UserType != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not admin"})
			return
		}

		cursor, err := userCollection.Find(c, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for cursor.Next(c) {
			var user models.AccountInfo
			cursor.Decode(&user)
			users = append(users, user)
		}

		c.JSON(http.StatusOK, gin.H{"result": users})
	}
}
