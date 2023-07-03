package main

import (
	"errors"
	"quickpolls/database"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var secret = []byte("secret")

type Claims struct {
	jwt.Claims
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       uint   `json:"id"`
}

func (c Claims) Valid() error {
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPasswordHash(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func CreateToken(user *database.User) (string, error) {
	claims := &Claims{
		Username: user.Username,
		Email:    user.Email,
		ID:       user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	if strings.Count(tokenString, ".") < 2 {
		return nil, errors.New("cadena de token no válida")
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if token.Valid && err == nil {
		return claims, nil
	}

	return nil, err
}
