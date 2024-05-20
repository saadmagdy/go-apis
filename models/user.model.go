package models

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID   `json:"_id" bson:"_id"`
	Name      string               `json:"name" bson:"name" validate:"required"`
	Email     string               `json:"email" bson:"email" validate:"required,email"`
	User_Type string               `json:"user_type" bson:"user_type" validate:"required,eq=SELLER|eq=BUYER"`
	Password  string               `json:"password,omitempty" bson:"password" validate:"required"`
	Products  []primitive.ObjectID `json:"products,omitempty" bson:"products,omitempty"`
}

type UserSignUp struct {
	Name      string `json:"name" bson:"name" validate:"required"`
	Email     string `json:"email" bson:"email" validate:"required,email"`
	User_Type string `json:"user_type" bson:"user_type" validate:"required,eq=SELLER|eq=BUYER"`
	Password  string `json:"password" bson:"password" validate:"required"`
}

type UserLogin struct {
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required"`
}

type UserUpdate struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email" validate:"email"`
	Password string `json:"password" bson:"password"`
}

func (u *User) GenerateToken() (string, error) {
	claims := jwt.MapClaims{
		"userId":    u.ID,
		"userEmail": u.Email,
		"userType":  u.User_Type,
		"exp":       time.Now().Add(time.Hour * 1).Unix(),
	}

	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwt_token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return token, nil
}
