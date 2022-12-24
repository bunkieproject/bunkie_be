package main

import (
	routes "github.com/bunkieproject/bunkie_be/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := "8080"

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.AdRoutes(router)

	router.Run(":" + port)
}
