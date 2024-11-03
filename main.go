package main

import (
	"net/http"
	"os"

	routes "github.com/bhaveshs012/golang-jwt-project/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger()) //* similar to adding middlewares / extra functionality to router

	// adding the router to the routes
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	//* API Checkers ::
	router.GET("/api/v1", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Success Access Granted for API V1",
		})
	})

	router.GET("/api/v2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Success Access Granted for API V2",
		})
	})
	router.Run(":" + port)
}
