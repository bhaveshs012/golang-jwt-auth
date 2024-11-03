package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bhaveshs012/golang-jwt-project/database"
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type signedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UserId    string
	UserType  string
	jwt.RegisteredClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &signedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UserId:    uid,
		UserType:  userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 24 hours expiry
		},
	}

	refreshClaims := &signedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 168)), // 7 days expiry
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}
func UpdateAllTokens(signedAccessToken, signedRefreshToken, userId string) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel() // Ensure the cancel function is called to prevent resource leaks

	// Prepare the update object with the $set operator
	updateObj := bson.D{
		{"$set", bson.D{
			{"access_token", signedAccessToken},
			{"refresh_token", signedRefreshToken},
			{"updated_at", time.Now()},
		}},
	}

	// Create filter for the update operation
	filter := bson.M{"user_id": userId}

	// Perform the update operation
	result, err := userCollection.UpdateOne(ctx, filter, updateObj)
	if err != nil {
		log.Printf("Error updating tokens: %v", err) // More informative error log
		return
	}

	// Log the result of the update operation
	if result.MatchedCount == 0 {
		log.Println("No documents matched the filter. Update not applied.")
	} else {
		log.Println("Tokens updated successfully.")
	}
}

func ValidateToken(signedToken string) (claims *signedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &signedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*signedDetails)
	if !ok {
		msg = fmt.Sprintf("The token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt.Before(time.Now()) {
		msg = "token is expired"
		return
	}

	return claims, msg
}
