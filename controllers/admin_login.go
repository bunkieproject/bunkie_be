package controllers

import (
	"net/http"

	"github.com/bunkieproject/bunkie_be/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func AdminLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.LoginRequest
		var user models.AccountInfo
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

		err := userCollection.FindOne(c, bson.M{"username": request.UsernameOrEmail}).Decode(&user)
		if err != nil { // if username not found, try email
			err = userCollection.FindOne(c, bson.M{"email": request.UsernameOrEmail}).Decode(&user)
			if err != nil { // if email not found, return error
				c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
				return
			}
		}

		isPassValid, err := VerifyPassword(*user.PasswordHash, *request.Password)
		if !isPassValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if *user.UserType != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You are not admin"})
			return
		}

		onlineToken.Token = user.Token
		onlineToken.User_id = user.User_id

		_, err = onlineCollection.InsertOne(c, onlineToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": onlineToken.Token})
	}
}
