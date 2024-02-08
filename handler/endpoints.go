package handler

import (
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// (POST /registration)
func (s *Server) Registration(ctx echo.Context) error {

	// Create RegistrationJSONRequestBody object
	var body generated.RegistrationJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid request body")
	}

	// Validate full name
	if len(body.FullName) < 3 || len(body.FullName) > 60 {
		return ctx.JSON(http.StatusBadRequest, "Full name must be between 3 and 60 characters")
	}

	// Validate phone number format
	phoneNumberRegex := regexp.MustCompile(`^\+62\d{8,11}$`)
	if !phoneNumberRegex.MatchString(body.PhoneNumber) {
		return ctx.JSON(http.StatusBadRequest, "Phone number must start with '+62' and have 10-13 digits")
	}

	// Validate password format
	if len(body.Password) < 6 || len(body.Password) > 64 {
		return ctx.JSON(http.StatusBadRequest, "Password must be between 6 and 64 characters")
	}
	if !ContainsCapital(body.Password) || !ContainsNumber(body.Password) || !ContainsSpecialCharacter(body.Password) {
		return ctx.JSON(http.StatusBadRequest, "Password must contain at least 1 capital letter, 1 number, and 1 special character")
	}

	// Hash and salt the password
	hashedPassword, err := HashPassword(body.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Failed to hash password")
	}

	input := repository.RegistrationInput{
		FullName:    body.FullName,
		PhoneNumber: body.PhoneNumber,
		Password:    hashedPassword,
	}

	// Call repository method for registration
	output, err := s.Repository.Registration(ctx.Request().Context(), input)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			// Handle duplicate key error
			return ctx.JSON(http.StatusBadRequest, "Phone number has been used to register another account")
		}
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, output)
}

// (POST /login)
func (s *Server) Login(ctx echo.Context) error {

	// Create LoginJSONRequestBody object
	var body generated.LoginJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid request body")
	}

	// Call repository method for login
	input := repository.LoginInput{
		PhoneNumber: body.PhoneNumber,
		Password:    body.Password,
	}
	output, err := s.Repository.Login(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid phone number")
	}

	// Verify the provided password against the hashed password stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(output.Password), []byte(body.Password)); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid password")
	}

	// Password is correct, proceed with login

	// Load RSA private key
	privateKeyBytes, err := os.ReadFile("private.pem")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	// Generate JWT token with RS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id": output.Id,
	})

	// Sign the token with RSA private key
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	// Create the response object containing the user ID and JWT token
	response := generated.LoginResponse{
		UserId: output.Id,
		Jwt:    tokenString,
	}

	// Return the response
	return ctx.JSON(http.StatusOK, response)
}

// (GET /my-profile)
func (s *Server) GetMyProfile(ctx echo.Context) error {
	// Get the JWT token from the Authorization header
	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ctx.JSON(http.StatusForbidden, "JWT token not found")
	}

	// Extract the token from the Authorization header
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ctx.JSON(http.StatusForbidden, "Invalid Authorization header format")
	}

	// Parse the token
	tokenString := parts[1]
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
		return ctx.JSON(http.StatusForbidden, err.Error())
	}

	// Check if the token is valid
	if !token.Valid {
		return ctx.JSON(http.StatusForbidden, "Invalid JWT token")
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ctx.JSON(http.StatusForbidden, "Failed to extract claims from JWT token")
	}

	// Extract the user ID from the claims
	userID, ok := claims["id"].(string)
	if !ok {
		return ctx.JSON(http.StatusForbidden, "Failed to extract user ID from JWT token")
	}

	// Use the user ID to retrieve profile information
	input := repository.GetMyProfileInput{
		Id: userID,
	}
	output, err := s.Repository.GetMyProfile(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, "Failed to retrieve user profile")
	}

	return ctx.JSON(http.StatusOK, output)
}

// (PUT /my-profile)
func (s *Server) UpdateMyProfile(ctx echo.Context) error {
	// Create UpdateMyProfileJSONRequestBody object
	var body generated.UpdateMyProfileJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusForbidden, "Invalid request body")
	}

	// Get the JWT token from the Authorization header
	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ctx.JSON(http.StatusForbidden, "JWT token not found")
	}

	// Extract the token from the Authorization header
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ctx.JSON(http.StatusForbidden, "Invalid Authorization header format")
	}

	// Parse the token
	tokenString := parts[1]
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
		return ctx.JSON(http.StatusForbidden, err.Error())
	}

	// Check if the token is valid
	if !token.Valid {
		return ctx.JSON(http.StatusForbidden, "Invalid JWT token")
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ctx.JSON(http.StatusForbidden, "Failed to extract claims from JWT token")
	}

	// Extract the user ID from the claims
	userID, ok := claims["id"].(string)
	if !ok {
		return ctx.JSON(http.StatusForbidden, "Failed to extract user ID from JWT token")
	}

	// Update the user profile
	input := repository.UpdateMyProfileInput{
		Id:          userID,
		FullName:    body.FullName,
		PhoneNumber: body.PhoneNumber,
	}
	output, err := s.Repository.UpdateMyProfile(ctx.Request().Context(), input)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			// Handle duplicate key error
			return ctx.JSON(http.StatusConflict, "Phone number has been used to register another account")
		}
		return ctx.JSON(http.StatusForbidden, "Failed to update user profile")
	}

	return ctx.JSON(http.StatusOK, output)
}
