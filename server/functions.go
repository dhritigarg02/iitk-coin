package server

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func CreateToken(rollno int) (string, error) {

	atClaims := jwt.MapClaims{
		"authorized": true,
		"rollno":     rollno,
		"exp":        time.Now().Add(time.Minute * 15).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims).SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
