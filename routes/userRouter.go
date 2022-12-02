package routes

import (
	controller "bunkie_be/controllers"
	"bunkie_be/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
}
