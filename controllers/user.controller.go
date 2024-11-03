package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bhaveshs012/golang-jwt-project/database"
	"github.com/bhaveshs012/golang-jwt-project/helpers"
	"github.com/bhaveshs012/golang-jwt-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var validate = validator.New()

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		context, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		var user models.User
		var foundUser models.User

		err := ctx.BindJSON(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = userCollection.FindOne(context, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User or password is incorrect"})
			return
		}

		passValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passValid {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.UserType, foundUser.UserId)
		helpers.UpdateAllTokens(token, refreshToken, foundUser.UserId)

		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.UserId}).Decode(&foundUser)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, foundUser)
	}
}

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		context, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		var user models.User //* golang model

		//* bind json to model -> basically jo request se aayega usko struct mein daalo
		err := ctx.BindJSON(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		//* Validate the struct
		validationErr := validate.Struct(user)
		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}

		//* check if this phone or email already exists : for mongo req always pass the timeout context and while returning res use the gin context
		count, err := userCollection.CountDocuments(context, bson.M{"email": user.Email})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		//* Hash the password
		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(context, bson.M{"phone_number": user.Phone})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error while Checking the phone number",
			})
			return
		}

		if count > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User with this email or phone number already exixts",
			})
			return
		}

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Id = primitive.NewObjectID()

		user.UserId = user.Id.Hex()

		//* Generate the tokens
		accessToken, refreshToken, err := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, user.UserId)
		if err != nil {
			msg := fmt.Sprintf("Error in creating tokens")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		//* setting the tokens
		user.AccessToken = &accessToken
		user.RefreshToken = &refreshToken

		resultInsertionNumber, insertError := userCollection.InsertOne(context, user)
		if insertError != nil {
			msg := fmt.Sprintf("User item was not created")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		ctx.JSON(http.StatusOK, resultInsertionNumber)

	}
}

func GetUserById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("user_id")

		err := helpers.MatchUserTypeToId(ctx, userId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		context, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		var user models.User

		//* Need to convert the mongo DB structure into struct for golang to understand
		err = userCollection.FindOne(context, bson.M{"_id": userId}).Decode(&user)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}
func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if the user type is ADMIN
		err := helpers.CheckUserType(ctx, "ADMIN")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create a context with a timeout
		contxt, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		// Get the record per page from the request, default to 10 if not provided or invalid
		recordPerPage, err := strconv.Atoi(ctx.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		// Get the page number, default to 1 if not provided or invalid
		page, err := strconv.Atoi(ctx.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		// Calculate the start index
		startIndex := (page - 1) * recordPerPage
		if queryStartIndex := ctx.Query("startIndex"); queryStartIndex != "" {
			startIndex, err = strconv.Atoi(queryStartIndex)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startIndex value"})
				return
			}
		}

		// Define aggregation stages
		matchStage := bson.D{{"$match", bson.D{}}}
		groupStage := bson.D{
			{"$group", bson.D{
				{"_id", nil},
				{"total_count", bson.D{{"$sum", 1}}},
				{"data", bson.D{{"$push", "$$ROOT"}}},
			}},
		}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
			}},
		}

		// Execute the aggregation pipeline
		cursor, err := userCollection.Aggregate(contxt, mongo.Pipeline{matchStage, groupStage, projectStage})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing user items"})
			return
		}

		// Read the aggregation results
		var allUsers []bson.M
		if err = cursor.All(contxt, &allUsers); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
			return
		}

		// Handle empty results
		if len(allUsers) == 0 {
			ctx.JSON(http.StatusOK, gin.H{"total_count": 0, "user_items": []bson.M{}})
			return
		}

		// Return the first result
		ctx.JSON(http.StatusOK, allUsers[0])
	}
}

func HashPassword(password string) string {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(pass)
}
func VerifyPassword(enteredPassword, actualPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(actualPassword), []byte(enteredPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Email or password is incorrect")
		check = false
	}
	return check, msg
}
