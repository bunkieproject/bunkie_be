package main

import (
	"log"
	"os"

	routes "github.com/bunkieproject/bunkie_be/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	port := "8080"

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.AdRoutes(router)

	router.Run(":" + port)
}
