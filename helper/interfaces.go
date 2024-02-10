package helper

import "github.com/dgrijalva/jwt-go"

type HelperInterface interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
	GenerateJWT(userId string) (string, error)
	ParseJWT(tokenString string) (jwt.MapClaims, error)
}
