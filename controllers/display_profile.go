package controllers

import (
	"net/http"

	"github.com/bunkieproject/bunkie_be/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
		if user.User_id != request.User_id { // other user's profile
			err = userCollection.FindOne(c, bson.M{"user_id": request.User_id}).Decode(&user)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
				return
			}

			if user.UserType != nil {
				if *user.UserType == "admin" {
					displayBanAndWarn = true
				} else {
					displayBanAndWarn = false
				}
			} else {
				displayBanAndWarn = false
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
				if user.ProfileInfo.DisplayEmail != nil {
					if !*user.ProfileInfo.DisplayEmail {
						user.Email = nil
					}
				} else {
					user.Email = nil
				}
				if user.ProfileInfo.DisplayPhone != nil {
					if !*user.ProfileInfo.DisplayPhone {
						user.ProfileInfo.Phone = nil
					}
				} else {
					user.ProfileInfo.Phone = nil
				}
			}

			if user.Email == nil {
				if user.ProfileInfo.Phone == nil {
					c.JSON(http.StatusOK, gin.H{"user_profile_info": user.ProfileInfo, "username": user.Username, "user_account_info": bson.M{"email": nil, "phone_number": nil}, "displayBanAndWarn": displayBanAndWarn, "room_ads": room_ads, "bunkie_ads": bunkie_ads})
					return
				} else {
					c.JSON(http.StatusOK, gin.H{"user_profile_info": user.ProfileInfo, "username": user.Username, "user_account_info": bson.M{"email": nil, "phone_number": *user.ProfileInfo.Phone}, "displayBanAndWarn": displayBanAndWarn, "room_ads": room_ads, "bunkie_ads": bunkie_ads})
					return
				}
			} else {
				if user.ProfileInfo.Phone == nil {
					c.JSON(http.StatusOK, gin.H{"user_profile_info": user.ProfileInfo, "username": user.Username, "user_account_info": bson.M{"email": *user.Email, "phone_number": nil}, "displayBanAndWarn": displayBanAndWarn, "room_ads": room_ads, "bunkie_ads": bunkie_ads})
					return
				} else {
					c.JSON(http.StatusOK, gin.H{"user_profile_info": user.ProfileInfo, "username": user.Username, "user_account_info": bson.M{"email": *user.Email, "phone_number": *user.ProfileInfo.Phone}, "displayBanAndWarn": displayBanAndWarn, "room_ads": room_ads, "bunkie_ads": bunkie_ads})
					return
				}
			}

		} else { // own profile
			displayBanAndWarn = false

			opts := options.Find()
			filter := bson.M{"user_id": request.User_id}
			cur, err := roomAdCollection.Find(c, filter, opts)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			for cur.Next(c) {
				var room_ad models.RoomAd
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

			for cur_bun.Next(c) {
				var bunkie_ad models.BunkieAd
				err := cur_bun.Decode(&bunkie_ad)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				bunkie_ads = append(bunkie_ads, bunkie_ad)
			}
			c.JSON(http.StatusOK, gin.H{"user_profile_info": user.ProfileInfo, "user_account_info": bson.M{"email": *user.Email, "username": *user.Username}, "displayBanAndWarn": displayBanAndWarn, "room_ads": room_ads, "bunkie_ads": bunkie_ads})

		}

	}
}
