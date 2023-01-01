package controllers

import (
	"net/http"

	"github.com/bunkieproject/bunkie_be/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func EditProfileInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.EditProfileInfoRequest
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

		err := userCollection.FindOne(c, bson.M{"token": request.Token}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.User_id != request.User_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User id does not match token"})
			return
		}

		err = userCollection.FindOne(c, bson.M{"user_id": request.User_id}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		if request.ProfileInfo.FullName != nil {
			user.ProfileInfo.FullName = request.ProfileInfo.FullName
		}
		if request.ProfileInfo.Phone != nil {
			user.ProfileInfo.Phone = request.ProfileInfo.Phone
		}
		if request.ProfileInfo.ProfilePicture != nil {
			user.ProfileInfo.ProfilePicture = request.ProfileInfo.ProfilePicture
		}
		if request.ProfileInfo.Description != nil {
			user.ProfileInfo.Description = request.ProfileInfo.Description
		}
		if request.ProfileInfo.Gender != nil {
			user.ProfileInfo.Gender = request.ProfileInfo.Gender
		}
		if request.ProfileInfo.City != nil {
			user.ProfileInfo.City = request.ProfileInfo.City
		}
		if request.ProfileInfo.DisplayEmail != nil {
			user.ProfileInfo.DisplayEmail = request.ProfileInfo.DisplayEmail
		}
		if request.ProfileInfo.DisplayPhone != nil {
			user.ProfileInfo.DisplayPhone = request.ProfileInfo.DisplayPhone
		}

		resultUpdateNumber, err := userCollection.UpdateOne(c, bson.M{"user_id": request.User_id}, bson.M{"$set": user})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": resultUpdateNumber})
	}

}
