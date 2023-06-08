package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/fahrizalfarid/user-service-grpc/src/constant"
	"github.com/fahrizalfarid/user-service-grpc/src/model"
	"github.com/fahrizalfarid/user-service-grpc/src/request"
	"github.com/fahrizalfarid/user-service-grpc/src/response"
)

type userSvcUsecase struct {
	User model.UserRepository
}

func NewUserSvcUsecase(user model.UserRepository) model.UserSvcUsecase {
	return &userSvcUsecase{
		User: user,
	}
}

func (u *userSvcUsecase) CreateUser(ctx context.Context, data *request.UserRequestService) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	rowsAffected, id, err := u.User.CreateUserProfile(ctx, &model.UserProfiles{
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Email:     data.Email,
		CreatedAt: data.CreatedAt,
		Phone:     data.Phone,
		Address:   data.Address,
		DeletedAt: int64(0),
	})
	if err != nil || rowsAffected == 0 {
		if strings.Contains(err.Error(), "unique constraint") {
			return 0, errors.New("email already exist")
		}
		return 0, err
	}

	rowsAffected, err = u.User.CreateUserCredential(ctx, &model.UserCredentials{
		Id:       id,
		Username: data.Username,
		Password: data.Password,
	})
	if err != nil || rowsAffected == 0 {
		if strings.Contains(err.Error(), "unique constraint") {
			return 0, errors.New("username already exist")
		}
		return 0, err
	}

	return id, nil
}
func (u *userSvcUsecase) GetUserById(ctx context.Context, id int64) (*response.UserProfileResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	user, err := u.User.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &response.UserProfileResponse{
		Id:        user.Id,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: time.Unix(user.CreatedAt, 0).Format(constant.TimeLayoutFormat),
		Phone:     user.Phone,
		Address:   user.Address,
		Username:  user.Username,
	}, nil
}

func (u *userSvcUsecase) GetUserCredentials(ctx context.Context, usernameOrEmail string) (*response.UserLoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	UserCredential, err := u.User.GetUserCredentials(ctx, usernameOrEmail)
	if err != nil {
		return nil, err
	}
	return UserCredential, nil
}

func (u *userSvcUsecase) GetUserByEmailOrUsername(ctx context.Context, username string) ([]*response.UserFound, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	user, err := u.User.GetUserByUsernameOrEmail(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userSvcUsecase) UpdateById(ctx context.Context, data *request.UserUpdateService) (*response.UserProfileResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rowsAffected, err := u.User.UpdateUserById(ctx, &model.UserProfiles{
		Id:        data.Id,
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Email:     data.Email,
		Phone:     data.Phone,
		Address:   data.Address,
	})
	if err != nil || rowsAffected == 0 {
		return nil, err
	}

	rowsAffected, err = u.User.UpdateUserCredentialById(ctx, &model.UserCredentials{
		Id:       data.Id,
		Username: data.Username,
		Password: data.Password,
	})

	if err != nil || rowsAffected == 0 {
		return nil, err
	}

	user, _ := u.User.GetUserById(ctx, data.Id)

	return &response.UserProfileResponse{
		Id:        user.Id,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: time.Unix(user.CreatedAt, 0).Format(constant.TimeLayoutFormat),
		Phone:     data.Phone,
		Address:   data.Address,
		Username:  data.Username,
	}, nil
}

func (u *userSvcUsecase) DeleteById(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	rowsAffected, err := u.User.DeleteUserById(ctx, id)
	if err != nil || rowsAffected == 0 {
		return err
	}
	return nil
}
