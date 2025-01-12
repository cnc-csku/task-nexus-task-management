package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID           bson.ObjectID `bson:"_id" json:"id"`
	Email        string        `bson:"email" json:"email"`
	PasswordHash string        `bson:"password_hash" json:"passwordHash"`
	FullName     string        `bson:"full_name" json:"fullName"`
	DisplayName  string        `bson:"display_name" json:"displayName"`
	CreatedAt    time.Time     `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time     `bson:"updated_at" json:"updatedAt"`
}

type UserCustomClaims struct {
	ID          string `json:"id"`
	FullName    string `json:"fullName"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	jwt.RegisteredClaims
}
