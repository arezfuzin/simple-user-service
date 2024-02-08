// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	Registration(ctx context.Context, input RegistrationInput) (output RegistrationOutput, err error)
	Login(ctx context.Context, input LoginInput) (output LoginOutput, err error)
	GetMyProfile(ctx context.Context, input GetMyProfileInput) (output GetMyProfileOutput, err error)
	UpdateMyProfile(ctx context.Context, input UpdateMyProfileInput) (output UpdateMyProfileOutput, err error)
}
