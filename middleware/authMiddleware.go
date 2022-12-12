package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error", "No Authorization token provided"})
			c.Abort()
			return
		}

		claims, msg := helper.ValidateToken(clientToken)
		if msg != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error", msg})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("username", claims.Username)
		c.Set("uid", claims.Uid)

		c.Next()
