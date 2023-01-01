package controllers

import (
	"net/http"

	"github.com/bunkieproject/bunkie_be/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func EnterNewPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.EnterNewPasswordRequest
		var user models.AccountInfo

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := userCollection.FindOne(c, bson.M{"email": request.Email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		password, _ := HashPassword(*request.NewPassword)
		err = userCollection.FindOneAndUpdate(c, bson.M{"email": request.Email}, bson.M{"$set": bson.M{"passwordhash": password, "password": *request.NewPassword}}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})

	}
}
