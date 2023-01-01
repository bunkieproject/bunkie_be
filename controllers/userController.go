package controllers

import (
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"time"

	"github.com/bunkieproject/bunkie_be/database"
	"github.com/bunkieproject/bunkie_be/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func checkIfUserOnline(user_id string, c *gin.Context) bool {
	var user models.AccountInfo

	err := onlineCollection.FindOne(c, bson.M{"user_id": user_id}).Decode(&user)
	if err != nil {
		return false
	}
	return true
}
