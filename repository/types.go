// This file contains types that are used in the repository layer.
package repository

// registration type
type RegistrationInput struct {
	FullName    string
	PhoneNumber string
	Password    string
}

type RegistrationOutput struct {
	UserId string
}

// login type
type LoginInput struct {
	PhoneNumber string
	Password    string
}

type LoginOutput struct {
	Id          string
	FullName    string
	PhoneNumber string
	Password    string
}

// get my profile type
type GetMyProfileInput struct {
	Id string
}

type GetMyProfileOutput struct {
	FullName    string
	PhoneNumber string
}

// update my profile type
type UpdateMyProfileInput struct {
	Id          string
	FullName    *string
	PhoneNumber *string
}

type UpdateMyProfileOutput struct {
	FullName    string
	PhoneNumber string
}
