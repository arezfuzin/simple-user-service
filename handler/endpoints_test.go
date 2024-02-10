package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/UserService/helper"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

type setupMock func(repo *repository.MockRepositoryInterface, helper *helper.MockHelperInterface)

func TestRegistration(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  map[string]interface{}
		setupMock    setupMock
		responseBody string
		statusCode   int
	}{
		{
			name:         "invalid full name",
			requestBody:  map[string]interface{}{},
			setupMock:    nil,
			responseBody: "\"Full name must be between 3 and 60 characters\"\n",
			statusCode:   http.StatusBadRequest,
		},
		{
			name: "phone number must start with '+62'",
			requestBody: map[string]interface{}{
				"fullName":    "John Doe",
				"phoneNumber": "081234567890",
				"password":    "Test123!",
			},
			setupMock:    nil,
			responseBody: "\"Phone number must start with '+62' and have 10-13 digits\"\n",
			statusCode:   http.StatusBadRequest,
		},
		{
			name: "phone number must have 10-13 digits",
			requestBody: map[string]interface{}{
				"fullName":    "John Doe",
				"phoneNumber": "+628123",
				"password":    "Test123!",
			},
			setupMock:    nil,
			responseBody: "\"Phone number must start with '+62' and have 10-13 digits\"\n",
			statusCode:   http.StatusBadRequest,
		},
		{
			name: "password must be between 6 and 64 characters",
			requestBody: map[string]interface{}{
				"fullName":    "John Doe",
				"phoneNumber": "+6281234567890",
				"password":    "ok!",
			},
			setupMock:    nil,
			responseBody: "\"Password must be between 6 and 64 characters\"\n",
			statusCode:   http.StatusBadRequest,
		},
		{
			name: "password must contain at least 1 capital letter, 1 number, and 1 special character",
			requestBody: map[string]interface{}{
				"fullName":    "John Doe",
				"phoneNumber": "+6281234567890",
				"password":    "aaaaaaaaaaaaaaaaaaaaaaa",
			},
			setupMock:    nil,
			responseBody: "\"Password must contain at least 1 capital letter, 1 number, and 1 special character\"\n",
			statusCode:   http.StatusBadRequest,
		},
		{
			name: "duplicate phone number",
			requestBody: map[string]interface{}{
				"fullName":    "John Doe",
				"phoneNumber": "+6281234567890",
				"password":    "Test123!",
			},
			setupMock: func(repo *repository.MockRepositoryInterface, helper *helper.MockHelperInterface) {
				helper.EXPECT().HashPassword(gomock.Any()).AnyTimes().Return("", nil)
				repo.EXPECT().Registration(gomock.Any(), gomock.Any()).AnyTimes().Return(repository.RegistrationOutput{}, &pq.Error{Code: "23505"})
			},
			responseBody: "\"Phone number has been used to register another account\"\n",
			statusCode:   http.StatusBadRequest,
		},
		{
			name: "failed to record registration",
			requestBody: map[string]interface{}{
				"fullName":    "John Doe",
				"phoneNumber": "+6281234567890",
				"password":    "Test123!",
			},
			setupMock: func(repo *repository.MockRepositoryInterface, helper *helper.MockHelperInterface) {
				helper.EXPECT().HashPassword(gomock.Any()).AnyTimes().Return("", nil)
				repo.EXPECT().Registration(gomock.Any(), gomock.Any()).AnyTimes().Return(repository.RegistrationOutput{}, errors.New("Failed to record registration"))
			},
			responseBody: "{}\n",
			statusCode:   http.StatusBadRequest,
		},
		{
			name: "failed to hash password",
			requestBody: map[string]interface{}{
				"fullName":    "John Doe",
				"phoneNumber": "+6281234567890",
				"password":    "Test123!",
			},
			setupMock: func(repo *repository.MockRepositoryInterface, helper *helper.MockHelperInterface) {
				helper.EXPECT().HashPassword(gomock.Any()).AnyTimes().Return("", errors.New("Failed to hash password"))
			},
			responseBody: "\"Failed to hash password\"\n",
			statusCode:   http.StatusBadRequest,
		},
		{
			name: "success",
			requestBody: map[string]interface{}{
				"fullName":    "John Doe",
				"phoneNumber": "+6281234567890",
				"password":    "Test123!",
			},
			setupMock: func(repo *repository.MockRepositoryInterface, helper *helper.MockHelperInterface) {
				helper.EXPECT().HashPassword(gomock.Any()).AnyTimes().Return("", nil)
				repo.EXPECT().Registration(gomock.Any(), gomock.Any()).AnyTimes().Return(repository.RegistrationOutput{
					UserId: "b3299b0b-2ffb-42e9-9a7b-daa1b0d10bbc",
				}, nil)
			},
			responseBody: "{\"UserId\":\"b3299b0b-2ffb-42e9-9a7b-daa1b0d10bbc\"}\n",
			statusCode:   http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			mockRepo := repository.NewMockRepositoryInterface(mockCtl)
			mockHelper := helper.NewMockHelperInterface(mockCtl)
			if tt.setupMock != nil {
				tt.setupMock(mockRepo, mockHelper)
			}

			// Create a new instance of the Server
			s := &Server{
				Repository: mockRepo,
				Helper:     mockHelper,
			}

			// Create a new Echo context for the test
			e := echo.New()

			jsonRequestBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/registration", bytes.NewBuffer(jsonRequestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Call the Registration handler function
			err := s.Registration(c)

			// Assert that no error occurred
			assert.NoError(t, err)

			// Assert the response status code
			assert.Equal(t, tt.statusCode, rec.Code)

			// Assert the response body
			assert.Equal(t, tt.responseBody, rec.Body.String())
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  map[string]interface{}
		setupMock    setupMock
		responseBody string
		statusCode   int
	}{
		{
			name: "invalid phone number",
			requestBody: map[string]interface{}{
				"phoneNumber": "xxxx",
				"password":    "Test123!",
			},
			setupMock:    nil,
			responseBody: "\"Phone number must start with '+62' and have 10-13 digits\"\n",
			statusCode:   http.StatusBadRequest,
		},
		{
			name: "failed to retrieve user",
			requestBody: map[string]interface{}{
				"phoneNumber": "+6281234567899",
				"password":    "Test123!",
			},
			setupMock: func(repo *repository.MockRepositoryInterface, helper *helper.MockHelperInterface) {
				repo.EXPECT().Login(gomock.Any(), gomock.Any()).AnyTimes().Return(repository.LoginOutput{}, errors.New("Invalid phone number"))
			},
			responseBody: "\"Invalid phone number\"\n",
			statusCode:   http.StatusBadRequest,
		},
		{
			name: "failed to generate JWT token",
			requestBody: map[string]interface{}{
				"phoneNumber": "+6281234567891",
				"password":    "Test123!",
			},
			setupMock: func(repo *repository.MockRepositoryInterface, helper *helper.MockHelperInterface) {
				helper.EXPECT().ComparePassword(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
				helper.EXPECT().GenerateJWT(gomock.Any()).AnyTimes().Return("", errors.New("Failed to generate JWT token"))
				repo.EXPECT().Login(gomock.Any(), gomock.Any()).AnyTimes().Return(repository.LoginOutput{
					Id:          "b3299b0b-2ffb-42e9-9a7b-daa1b0d10bbc",
					FullName:    "John Doe",
					PhoneNumber: "+6281234567890",
					Password:    "Test123!",
				}, nil)
			},
			responseBody: "\"Failed to generate JWT token\"\n",
			statusCode:   http.StatusBadRequest,
		},
		{
			name: "success",
			requestBody: map[string]interface{}{
				"phoneNumber": "+6281234567890",
				"password":    "Test123!",
			},
			setupMock: func(repo *repository.MockRepositoryInterface, helper *helper.MockHelperInterface) {
				helper.EXPECT().ComparePassword(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
				helper.EXPECT().GenerateJWT(gomock.Any()).AnyTimes().Return("jwtToken", nil)
				repo.EXPECT().Login(gomock.Any(), gomock.Any()).AnyTimes().Return(repository.LoginOutput{
					Id:          "b3299b0b-2ffb-42e9-9a7b-daa1b0d10bbc",
					FullName:    "John Doe",
					PhoneNumber: "+6281234567890",
					Password:    "Test123!",
				}, nil)
			},
			responseBody: "{\"jwt\":\"jwtToken\",\"userId\":\"b3299b0b-2ffb-42e9-9a7b-daa1b0d10bbc\"}\n",
			statusCode:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			mockRepo := repository.NewMockRepositoryInterface(mockCtl)
			mockHelper := helper.NewMockHelperInterface(mockCtl)
			if tt.setupMock != nil {
				tt.setupMock(mockRepo, mockHelper)
			}

			// Create a new instance of the Server
			s := &Server{
				Repository: mockRepo,
				Helper:     mockHelper,
			}

			// Create a new Echo context for the test
			e := echo.New()

			jsonRequestBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonRequestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Call the Registration handler function
			err := s.Login(c)

			// Assert that no error occurred
			assert.NoError(t, err)

			// Assert the response status code
			assert.Equal(t, tt.statusCode, rec.Code)

			// Assert the response body
			assert.Equal(t, tt.responseBody, rec.Body.String())
		})
	}
}

func TestGetMyProfile(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        map[string]interface{}
		requestTokenHeader string
		setupMock          setupMock
		responseBody       string
		statusCode         int
	}{
		{
			name:               "token not found",
			requestBody:        nil,
			requestTokenHeader: "",
			setupMock:          nil,
			responseBody:       "\"JWT token not found\"\n",
			statusCode:         http.StatusForbidden,
		},
		{
			name:               "token invalid format",
			requestBody:        nil,
			requestTokenHeader: "Bearerinvalid_jwt_token",
			setupMock:          nil,
			responseBody:       "\"Invalid Authorization header format\"\n",
			statusCode:         http.StatusForbidden,
		},
		{
			name:               "failed to parse JWT token",
			requestBody:        nil,
			requestTokenHeader: "Bearer valid_jwt_token",
			setupMock: func(repo *repository.MockRepositoryInterface, helper *helper.MockHelperInterface) {
				helper.EXPECT().ParseJWT(gomock.Any()).AnyTimes().Return(map[string]interface{}{}, errors.New("Failed to parse JWT token"))
			},
			responseBody: "\"Failed to parse JWT token\"\n",
			statusCode:   http.StatusForbidden,
		},
		{
			name:               "failed to claim user ID from JWT token",
			requestBody:        nil,
			requestTokenHeader: "Bearer valid_jwt_token",
			setupMock: func(repo *repository.MockRepositoryInterface, helper *helper.MockHelperInterface) {
				helper.EXPECT().ParseJWT(gomock.Any()).AnyTimes().Return(map[string]interface{}{
					"userId": "b3299b0b-2ffb-42e9-9a7b-daa1b0d10bbc",
				}, nil)
			},
			responseBody: "\"Failed to extract user ID from JWT token\"\n",
			statusCode:   http.StatusForbidden,
		},
		{
			name:               "failed retrieve user profile",
			requestBody:        nil,
			requestTokenHeader: "Bearer valid_jwt_token",
			setupMock: func(repo *repository.MockRepositoryInterface, helper *helper.MockHelperInterface) {
				helper.EXPECT().ParseJWT(gomock.Any()).AnyTimes().Return(map[string]interface{}{
					"id": "b3299b0b-2ffb-42e9-9a7b-daa1b0d10bbc",
				}, nil)
				repo.EXPECT().GetMyProfile(gomock.Any(), gomock.Any()).AnyTimes().Return(repository.GetMyProfileOutput{}, errors.New("Failed to retrieve user profile"))
			},
			responseBody: "\"Failed to retrieve user profile\"\n",
			statusCode:   http.StatusForbidden,
		},
		{
			name:               "success",
			requestBody:        nil,
			requestTokenHeader: "Bearer valid_jwt_token",
			setupMock: func(repo *repository.MockRepositoryInterface, helper *helper.MockHelperInterface) {
				helper.EXPECT().ParseJWT(gomock.Any()).AnyTimes().Return(map[string]interface{}{
					"id": "b3299b0b-2ffb-42e9-9a7b-daa1b0d10bbc",
				}, nil)
				repo.EXPECT().GetMyProfile(gomock.Any(), gomock.Any()).AnyTimes().Return(repository.GetMyProfileOutput{
					FullName:    "John Doe",
					PhoneNumber: "+6281234567890",
				}, nil)
			},
			responseBody: "{\"FullName\":\"John Doe\",\"PhoneNumber\":\"+6281234567890\"}\n",
			statusCode:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			mockRepo := repository.NewMockRepositoryInterface(mockCtl)
			mockHelper := helper.NewMockHelperInterface(mockCtl)
			if tt.setupMock != nil {
				tt.setupMock(mockRepo, mockHelper)
			}

			// Create a new instance of the Server
			s := &Server{
				Repository: mockRepo,
				Helper:     mockHelper,
			}

			// Create a new Echo context for the test
			e := echo.New()

			req := httptest.NewRequest(http.MethodPost, "/my-profile", nil)
			req.Header.Set("Authorization", tt.requestTokenHeader)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Call the Registration handler function
			err := s.GetMyProfile(c)

			// Assert that no error occurred
			assert.NoError(t, err)

			// Assert the response status code
			assert.Equal(t, tt.statusCode, rec.Code)

			// Assert the response body
			assert.Equal(t, tt.responseBody, rec.Body.String())
		})
	}
}

func TestUpdateMyProfile(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        map[string]interface{}
		requestTokenHeader string
		setupMock          setupMock
		responseBody       string
		statusCode         int
	}{
		{
			name:               "token not found",
			requestBody:        map[string]interface{}{"fullName": "John Doe", "phoneNumber": "+6281234567890"},
			requestTokenHeader: "",
			setupMock:          nil,
			responseBody:       "\"JWT token not found\"\n",
			statusCode:         http.StatusForbidden,
		},
		{
			name:               "invalid format",
			requestBody:        map[string]interface{}{"fullName": "John Doe", "phoneNumber": "+6281234567890"},
			requestTokenHeader: "Bearerinvalid_jwt_token",
			setupMock:          nil,
			responseBody:       "\"Invalid Authorization header format\"\n",
			statusCode:         http.StatusForbidden,
		},
		{
			name:               "failed to parse JWT token",
			requestBody:        map[string]interface{}{"fullName": "John Doe", "phoneNumber": "+6281234567890"},
			requestTokenHeader: "Bearer valid_jwt_token",
			setupMock: func(repo *repository.MockRepositoryInterface, helper *helper.MockHelperInterface) {
				helper.EXPECT().ParseJWT(gomock.Any()).AnyTimes().Return(map[string]interface{}{}, errors.New("Failed to parse JWT token"))
			},
			responseBody: "\"Failed to parse JWT token\"\n",
			statusCode:   http.StatusForbidden,
		},
		{
			name:               "failed to claim user ID from JWT token",
			requestBody:        map[string]interface{}{"fullName": "John Doe", "phoneNumber": "+6281234567890"},
			requestTokenHeader: "Bearer valid_jwt_token",
			setupMock: func(repo *repository.MockRepositoryInterface, helper *helper.MockHelperInterface) {
				helper.EXPECT().ParseJWT(gomock.Any()).AnyTimes().Return(map[string]interface{}{
					"userId": "b3299b0b-2ffb-42e9-9a7b-daa1b0d10bbc",
				}, nil)
			},
			responseBody: "\"Failed to extract user ID from JWT token\"\n",
			statusCode:   http.StatusForbidden,
		},
		{
			name:               "failed to update user profile",
			requestBody:        map[string]interface{}{"fullName": "John Doe", "phoneNumber": "+6281234567890"},
			requestTokenHeader: "Bearer valid_jwt_token",
			setupMock: func(repo *repository.MockRepositoryInterface, helper *helper.MockHelperInterface) {
				helper.EXPECT().ParseJWT(gomock.Any()).AnyTimes().Return(map[string]interface{}{
					"id": "b3299b0b-2ffb-42e9-9a7b-daa1b0d10bbc",
				}, nil)
				repo.EXPECT().UpdateMyProfile(gomock.Any(), gomock.Any()).AnyTimes().Return(repository.UpdateMyProfileOutput{}, errors.New("Failed to update user profile"))
			},
			responseBody: "\"Failed to update user profile\"\n",
			statusCode:   http.StatusForbidden,
		},
		{
			name:               "success",
			requestBody:        map[string]interface{}{"fullName": "John Doe", "phoneNumber": "+6281234567890"},
			requestTokenHeader: "Bearer valid_jwt_token",
			setupMock: func(repo *repository.MockRepositoryInterface, helper *helper.MockHelperInterface) {
				helper.EXPECT().ParseJWT(gomock.Any()).AnyTimes().Return(map[string]interface{}{
					"id": "b3299b0b-2ffb-42e9-9a7b-daa1b0d10bbc",
				}, nil)
				repo.EXPECT().UpdateMyProfile(gomock.Any(), gomock.Any()).AnyTimes().Return(repository.UpdateMyProfileOutput{
					FullName:    "John Doe",
					PhoneNumber: "+6281234567890",
				}, nil)
			},
			responseBody: "{\"FullName\":\"John Doe\",\"PhoneNumber\":\"+6281234567890\"}\n",
			statusCode:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			mockRepo := repository.NewMockRepositoryInterface(mockCtl)
			mockHelper := helper.NewMockHelperInterface(mockCtl)
			if tt.setupMock != nil {
				tt.setupMock(mockRepo, mockHelper)
			}

			// Create a new instance of the Server
			s := &Server{
				Repository: mockRepo,
				Helper:     mockHelper,
			}

			// Create a new Echo context for the test
			e := echo.New()

			jsonRequestBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/my-profile", bytes.NewBuffer(jsonRequestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set("Authorization", tt.requestTokenHeader)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Call the Registration handler function
			err := s.UpdateMyProfile(c)

			// Assert that no error occurred
			assert.NoError(t, err)

			// Assert the response status code
			assert.Equal(t, tt.statusCode, rec.Code)

			// Assert the response body
			assert.Equal(t, tt.responseBody, rec.Body.String())
		})
	}
}
