package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Rollno    int       `json:"rollno"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(rollno int) (*Payload, error) {

	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("[NewPayload] : %v", err)
	}

	payload := &Payload{
		ID:        tokenID,
		Rollno:    rollno,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Minute * 15),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func CreateToken(rollno int) (string, error) {

	payload, err := NewPayload(rollno)
	if err != nil {
		return "", fmt.Errorf("[CreateToken] : %v", err)
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", fmt.Errorf("[CreateToken] : %v", err)
	}
	return token, nil
}

func VerifyToken(tokenString string) (*Payload, error) {

	keyfunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, keyfunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
