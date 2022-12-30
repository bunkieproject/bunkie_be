package controllers

import (
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/bunkieproject/bunkie_be/database"
	helper "github.com/bunkieproject/bunkie_be/helpers"
	"github.com/bunkieproject/bunkie_be/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var onlineCollection *mongo.Collection = database.OpenCollection(database.Client, "online")
var bannedUsersCollection *mongo.Collection = database.OpenCollection(database.Client, "banned_users")
var warnedUsersCollection *mongo.Collection = database.OpenCollection(database.Client, "warned_users")
var validate = validator.New()

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

// VerifyPassword compares the password with the hash
func VerifyPassword(userPassword string, providedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	if err != nil {
		return false, err
	}
	return true, nil
}

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

		resultInsertNumber, err := userCollection.InsertOne(c, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": resultInsertNumber})
	}
}

// Login logs in a user
func Login() gin.HandlerFunc {
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

		err = bannedUsersCollection.FindOne(c, bson.M{"user_id": user.User_id}).Decode(&user)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You are banned"})
			return
		}

		err = warnedUsersCollection.FindOne(c, bson.M{"user_id": user.User_id}).Decode(&user)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You are warned"})
			return
		}

		onlineToken.Token = user.Token
		onlineToken.User_id = user.User_id

		_, err = onlineCollection.InsertOne(c, onlineToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": onlineToken.Token, "user_id": onlineToken.User_id})
	}
}

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

func ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.ResetPasswordRequest
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

		// generate new password
		verificationCode := GenerateSixDigit()

		// send this 6 digit verficiation code to user's email
		err = SendCodeToEmail(request.Email, verificationCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"verification_code": verificationCode})
	}
}

func GenerateSixDigit() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}

func SendCodeToEmail(email string, code string) error {
	from := os.Getenv("GOOGLE_EMAIL")
	to := email
	password := os.Getenv("GOOGLE_PASSWORD")
	msg := "From: " + from + " \n" + "To: " + to + " \n" + "Subject: Verification Code \n\n" + "Your verification code is: " + code
	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, password, "smtp.gmail.com"), from, []string{to}, []byte(msg))
	if err != nil {
		log.Fatal(err)
	}

	return err
}

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

func DeleteAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.DeleteAccountRequest
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

		if user.User_id != request.User_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		// delete user from database
		resultDeleteUser, err := userCollection.DeleteOne(c, bson.M{"token": request.Token})
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

		c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully", "resultDeleteUser": resultDeleteUser, "resultDeletePostsBunkie": resultDeletePostsBunkie, "resultDeletePostsRoom": resultDeletePostsRoom})
	}
}

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

func GetProfileInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id models.GetProfileInfoRequest
		var user models.AccountInfo

		if err := c.BindJSON(&id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(id)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := userCollection.FindOne(c, bson.M{"token": id.Token}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.User_id != id.User_id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User id does not match token"})
			return
		}

		err = userCollection.FindOne(c, bson.M{"user_id": id.User_id}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"profile_info": user.ProfileInfo})
	}
}

func UpdateAccountInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.UpdateAccountInfoRequest
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

		if request.Username != nil {
			user.Username = request.Username
		}
		if request.Email != nil {
			user.Email = request.Email
		}
		if request.Password != nil {
			user.Password = request.Password
		}

		resultUpdateNumber, err := userCollection.UpdateOne(c, bson.M{"user_id": request.User_id}, bson.M{"$set": user})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": resultUpdateNumber})
	}
}

func CreateProfileInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.CreateProfileInfoRequest
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

		user.ProfileInfo = request.ProfileInfo

		resultUpdateNumber, err := userCollection.UpdateOne(c, bson.M{"user_id": request.User_id}, bson.M{"$set": user})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": resultUpdateNumber})

	}

}

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

func DisplayProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.DisplayProfileRequest
		var user models.AccountInfo
		var room_ads []models.RoomAd
		var bunkie_ads []models.BunkieAd

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

		var displayBanAndWarn bool
		if *user.UserType == "admin" {
			displayBanAndWarn = true
		}

		err = userCollection.FindOne(c, bson.M{"user_id": request.User_id}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		opts := options.Find()
		filter := bson.M{"user_id": request.User_id}
		cur, err := roomAdCollection.Find(c, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var room_ad models.RoomAd
		for cur.Next(c) {
			err := cur.Decode(&room_ad)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			room_ads = append(room_ads, room_ad)
		}

		cur_bun, err := bunkieAdCollection.Find(c, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var bunkie_ad models.BunkieAd
		for cur_bun.Next(c) {
			err := cur_bun.Decode(&bunkie_ad)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			bunkie_ads = append(bunkie_ads, bunkie_ad)
		}

		if user.ProfileInfo != nil {
			if !*user.ProfileInfo.DisplayEmail {
				user.Email = nil
			}
			if !*user.ProfileInfo.DisplayPhone {
				user.ProfileInfo.Phone = nil
			}
		}
		c.JSON(http.StatusOK, gin.H{"user": user, "displayBanAndWarn": displayBanAndWarn, "room_ads": room_ads, "bunkie_ads": bunkie_ads})
	}
}

func BanUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.BanUserRequest
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if *user.UserType != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not admin"})
			return
		}

		err = userCollection.FindOne(c, bson.M{"user_id": request.User_id}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		err = onlineCollection.FindOne(c, bson.M{"user_id": request.User_id}).Decode(&user)
		if err == nil {
			_, err = onlineCollection.DeleteOne(c, bson.M{"user_id": request.User_id})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		_, err = bannedUsersCollection.InsertOne(c, bson.M{"user_id": request.User_id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": "User banned"})
	}
}

func WarnUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.WarnUserRequest
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if *user.UserType != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not admin"})
			return
		}

		err = userCollection.FindOne(c, bson.M{"user_id": request.User_id}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		_, err = warnedUsersCollection.InsertOne(c, bson.M{"user_id": request.User_id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": "User warned"})
	}
}

func UnbanUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.UnbanUserRequest
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if *user.UserType != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is not admin"})
			return
		}

		err = userCollection.FindOne(c, bson.M{"user_id": request.User_id}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		_, err = bannedUsersCollection.DeleteOne(c, bson.M{"user_id": request.User_id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": "User unbaned"})
	}
}

func checkIfUserOnline(user_id string, c *gin.Context) bool {
	var user models.AccountInfo

	err := onlineCollection.FindOne(c, bson.M{"user_id": user_id}).Decode(&user)
	if err != nil {
		return false
	}
	return true
}

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
