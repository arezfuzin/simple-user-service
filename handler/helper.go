package handler

import (
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
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
