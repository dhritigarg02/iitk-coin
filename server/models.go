package server

import (
	"time"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type User struct {
	Name     string `json:"name"`
	RollNo   int    `json:"rollno"`
	Password string `json:"password"`
	Batch    int    `json:"batch"`
}

type AuthUser struct {
	RollNo   int    `json:"rollno"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Rollno          int `json:"rollno"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}
