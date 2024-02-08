package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegistration(t *testing.T) {
	// Create a new instance of the Server
	s := &Server{}

	// Create a new Echo context for the test
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/registration", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the Registration handler function
	err := s.Registration(c)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Assert the response body
	assert.Equal(t, "Registration successful", rec.Body.String())
}

func TestLogin(t *testing.T) {
	// Create a new instance of the Server
	s := &Server{}

	// Create a new Echo context for the test
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the Login handler function
	err := s.Login(c)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Assert the response body
	assert.Equal(t, "Login successful", rec.Body.String())
}

func TestGetMyProfile(t *testing.T) {
	// Create a new instance of the Server
	s := &Server{}

	// Create a new Echo context for the test
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/my-profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the GetMyProfile handler function
	err := s.GetMyProfile(c)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Assert the response body
	assert.Equal(t, "My profile retrieved successfully", rec.Body.String())
}

func TestUpdateMyProfile(t *testing.T) {
	// Create a new instance of the Server
	s := &Server{}

	// Create a new Echo context for the test
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/my-profile", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the UpdateMyProfile handler function
	err := s.UpdateMyProfile(c)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Assert the response body
	assert.Equal(t, "My profile updated successfully", rec.Body.String())
}
