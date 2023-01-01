package controllers

import (
	"net/http"

	"github.com/bunkieproject/bunkie_be/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.DeleteUserRequest
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

		err := userCollection.FindOne(c, bson.M{"token": request.Token}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		if *user.UserType != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You are not admin"})
			return
		}

		err = userCollection.FindOne(c, bson.M{"user_id": request.User_id}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		// delete user from database
		resultDeleteUser, err := userCollection.DeleteOne(c, bson.M{"user_id": request.User_id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// delete user's posts from database
		resultDeletePostsBunkie, err := bunkieAdCollection.DeleteMany(c, bson.M{"user_id": request.User_id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resultDeletePostsRoom, err := roomAdCollection.DeleteMany(c, bson.M{"user_id": request.User_id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully", "resultDeleteUser": resultDeleteUser, "resultDeletePostsBunkie": resultDeletePostsBunkie, "resultDeletePostsRoom": resultDeletePostsRoom})
	}
}
