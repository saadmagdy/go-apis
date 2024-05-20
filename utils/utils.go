package utils

import (
	"context"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func Ctx() (context.Context, context.CancelFunc) {
	ctx, cansel := context.WithTimeout(context.Background(), 20*time.Second)
	return ctx, cansel
}

func HashPWD(password string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPwd), nil
}

func VerfiyPWD(givenPassword, userPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(givenPassword))
	msg := ""
	valid := true
	if err != nil {
		valid = false
		msg = "Invalid Password!"
	}
	return valid, msg
}

func VerifyToken(signdeToken string) (claims jwt.MapClaims, msg string) {
	token, err := jwt.Parse(signdeToken,  func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		msg = err.Error()
		return
	}

	claims ,ok:= token.Claims.(jwt.MapClaims)
	if !ok {
		msg = "the token is invalid"
		return
	}
	return claims, msg
}
