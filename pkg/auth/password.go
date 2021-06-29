package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPswd(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("[HashPswd] : %v", err)
	}
	return string(hash), err
}

func CheckPswd(password string, hashedpswd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedpswd), []byte(password))
}