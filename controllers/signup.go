package controllers

import (
	"net/http"
	"time"

	helper "github.com/bunkieproject/bunkie_be/helpers"
	"github.com/bunkieproject/bunkie_be/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SignUp creates a new user
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.SignUpRequest
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

		count_username, err := userCollection.CountDocuments(c, bson.M{"username": request.Username})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		count_email, err := userCollection.CountDocuments(c, bson.M{"email": request.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count_username > 0 || count_email > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}

		password, _ := HashPassword(*request.Password)

		user.Email = request.Email
		user.Username = request.Username
		user.Password = request.Password
		user.PasswordHash = &password
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		user.CreatedAt = time.Now()

		token, err := helper.GenerateToken(*user.Email, *user.Username, user.User_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		user.Token = &token
		user.ProfileInfo = &models.ProfileInfo{}

		resultInsertNumber, err := userCollection.InsertOne(c, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": resultInsertNumber})
	}
}
