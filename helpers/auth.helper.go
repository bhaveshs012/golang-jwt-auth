package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// * Getting the data stored in the context :: (need to save them in context first)
func CheckUserType(c *gin.Context, role string) error {
	userType := c.GetString("user_type")
	if userType != role {
		return errors.New("unauthorized to access this resource")
	}
	return nil
}

func MatchUserTypeToId(c *gin.Context, userId string) error {
	userType := c.GetString("user_type")
	uid := c.GetString("user_id")
	if userType == "USER" && uid != userId {
		return errors.New("unauthorized to access this resource")
	}
	return CheckUserType(c, userType)
}
