package usecase

import (
	"context"
	"time"

	"github.com/fahrizalfarid/user-service-grpc/src/model"
)

type userValidatorSvcUsecase struct {
	UserValidator model.ValidatorRepository
}

func NewUserValidatorSvcUsecase(userValidator model.ValidatorRepository) model.UserValidatorSvcUsecase {
	return &userValidatorSvcUsecase{
		UserValidator: userValidator,
	}
}

func (u *userValidatorSvcUsecase) IsUsernameExists(ctx context.Context, username string) bool {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	exists := u.UserValidator.IsUsernameExists(ctx, username)
	return exists
}

func (u *userValidatorSvcUsecase) IsEmailExists(ctx context.Context, email string) bool {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	exists := u.UserValidator.IsEmailExists(ctx, email)
	return exists
}

func (u *userValidatorSvcUsecase) IsUserExists(ctx context.Context, emailOrUsername string) bool {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	exists := u.UserValidator.IsUserExists(ctx, emailOrUsername)
	return exists
}
