package routes

import (
	"github.com/bhaveshs012/golang-jwt-project/controllers"
	"github.com/bhaveshs012/golang-jwt-project/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	// middleware injection
	incomingRoutes.Use(middlewares.Authenticate())

	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUserById())
}
