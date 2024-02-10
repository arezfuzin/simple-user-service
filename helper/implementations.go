package helper

import (
	"os"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Function to check if a string contains at least 1 capital letter
func ContainsCapital(s string) bool {
	for _, c := range s {
		if 'A' <= c && c <= 'Z' {
			return true
		}
	}
	return false
}

// Function to check if a string contains at least 1 number
func ContainsNumber(s string) bool {
	for _, c := range s {
		if '0' <= c && c <= '9' {
			return true
		}
	}
	return false
}

// Function to check if a string contains at least 1 special character
func ContainsSpecialCharacter(s string) bool {
	specialCharacters := "~!@#$%^&*()-_+=<>?/,.:;"
	for _, c := range s {
		for _, special := range specialCharacters {
			if c == special {
				return true
			}
		}
	}
	return false
}

// Function to hash and salt the password
func (*Helper) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Function to compare the password with the hashed password
func (*Helper) ComparePassword(hashedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err
	}

	return nil
}

// Function to generate JWT token
func (*Helper) GenerateJWT(userId string) (string, error) {
	// Load RSA private key
	privateKeyBytes, err := os.ReadFile("private.pem")
	if err != nil {
		return "", err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return "", err
	}

	// Generate JWT token with RS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id": userId,
	})

	// Sign the token with RSA private key
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Function to parse JWT token
func (*Helper) ParseJWT(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Load RSA private key
		privateKeyBytes, err := os.ReadFile("private.pem")
		if err != nil {
			return nil, err
		}

		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
		if err != nil {
			return nil, err
		}
		return privateKey.Public(), nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, err
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
