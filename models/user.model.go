package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id"` // Added `omitempty` for optional field in BSON
	FirstName    *string            `json:"first_name" validate:"required,min=2,max=100"`
	LastName     *string            `json:"last_name" validate:"required,min=2,max=100"`
	Email        *string            `json:"email" validate:"required,email"`
	Phone        *string            `json:"phone_number" validate:"required"`
	Password     *string            `json:"password" validate:"required,min=6"`
	AccessToken  *string            `json:"access_token,omitempty"`  // Added `omitempty` to skip if empty
	RefreshToken *string            `json:"refresh_token,omitempty"` // Added `omitempty` to skip if empty
	UserType     *string            `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	UserId       string             `json:"user_id" bson:"user_id"` // Added JSON and BSON tags for consistency
}
