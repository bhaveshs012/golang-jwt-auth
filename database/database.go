package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func connectDB() *mongo.Client {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading the dotenv File")
	}

	mongoDBUrl := os.Getenv("MONGODB_URL")
	clientOpts := options.Client().ApplyURI(mongoDBUrl)
	client, err := mongo.Connect(clientOpts)

	if err != nil {
		log.Fatal("Error connecting to the database")
	}
	fmt.Println("Connected to MongoDB !!")
	return client
}

var Client *mongo.Client = connectDB()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("golang-jwt").Collection(collectionName)
}
