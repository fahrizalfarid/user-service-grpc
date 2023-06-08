package model

import "context"

type (
	ValidatorRepository interface {
		IsUsernameExists(ctx context.Context, username string) bool
		IsEmailExists(ctx context.Context, email string) bool
		IsUserExists(ctx context.Context, usernameOrEmail string) bool
	}
	UserValidatorSvcUsecase interface {
		IsUsernameExists(ctx context.Context, username string) bool
		IsEmailExists(ctx context.Context, email string) bool
		IsUserExists(ctx context.Context, emailOrUsername string) bool
	}
)
