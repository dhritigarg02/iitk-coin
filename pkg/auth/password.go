package auth

import (

	"golang.org/x/crypto/bcrypt"
)

func HashPswd(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func CheckPswd(password string, hashedpswd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedpswd), []byte(password))
}