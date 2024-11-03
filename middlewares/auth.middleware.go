package middlewares

import (
	"net/http"

	"github.com/bhaveshs012/golang-jwt-project/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//* Get the token from the header :: and check if empty
		clientToken := ctx.Request.Header.Get("access_token")
		if clientToken == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization Token present"})
			ctx.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			ctx.Abort()
			return
		}
		ctx.Set("email", claims.Email)
		ctx.Set("first_name", claims.FirstName)
		ctx.Set("last_name", claims.LastName)
		ctx.Set("user_id", claims.UserId)
		ctx.Set("user_type", claims.UserType)
		ctx.Next()

	}
}
