package model

import (
	"context"

	pb "github.com/fahrizalfarid/user-service-grpc/src/proto"
	"github.com/fahrizalfarid/user-service-grpc/src/request"
	"github.com/fahrizalfarid/user-service-grpc/src/response"
)

type UserProfiles struct {
	Id        int64  `gorm:"primaryKey;type:int8;" json:"id"`
	Firstname string `gorm:"type:varchar(255);" json:"firstname"`
	Lastname  string `gorm:"type:varchar(255);" json:"lastname"`
	Email     string `gorm:"type:varchar(255);" json:"email"`
	CreatedAt int64  `gorm:"type:int8;" json:"created_at"`
	Phone     string `gorm:"type:varchar(20);" json:"phone"`
	Address   string `gorm:"type:varchar(255);" json:"address"`
	DeletedAt int64  `gorm:"type:int8;" json:"deleted_at"`
}

type UserCredentials struct {
	Id       int64  `gorm:"primaryKey;type:int64;" json:"id"`
	Username string `gorm:"type:varchar(255);" json:"username"`
	Password string `gorm:"type:varchar(255);"`
}

type (
	UserRepository interface {
		CreateUserProfile(ctx context.Context, data *UserProfiles) (int64, int64, error)
		CreateUserCredential(ctx context.Context, dataCreds *UserCredentials) (int64, error)
		GetUserById(ctx context.Context, id int64) (*response.UserProfile, error)
		GetUserCredentials(ctx context.Context, usernameOrEmail string) (*response.UserLoginResponse, error)
		GetUserByUsernameOrEmail(ctx context.Context, id string) ([]*response.UserFound, error)
		UpdateUserById(ctx context.Context, data *UserProfiles) (int64, error)
		UpdateUserCredentialById(ctx context.Context, data *UserCredentials) (int64, error)
		DeleteUserById(ctx context.Context, id int64) (int64, error)
	}
	UserSvcUsecase interface {
		CreateUser(ctx context.Context, data *request.UserRequestService) (int64, error)
		GetUserById(ctx context.Context, id int64) (*response.UserProfileResponse, error)
		GetUserCredentials(ctx context.Context, usernameOrEmail string) (*response.UserLoginResponse, error)
		GetUserByEmailOrUsername(ctx context.Context, username string) ([]*response.UserFound, error)
		UpdateById(ctx context.Context, data *request.UserUpdateService) (*response.UserProfileResponse, error)
		DeleteById(ctx context.Context, id int64) error
	}
	UserUsecase interface {
		CreateUser(ctx context.Context, data *request.UserRequest) (string, int64, error)
		GetUserById(ctx context.Context, id int64, token string) (*response.UserProfileResponse, error)
		Login(ctx context.Context, username, password string) (*response.AuthResponse, error)
		Find(ctx context.Context, username, token string) ([]*response.UserFound, error)
		FindWithArray(ctx context.Context, username, token string) (*pb.UserFoundArray, error)
		UpdateById(ctx context.Context, data *request.UserUpdateRequest, token string) (*response.UserProfileResponse, error)
		DeleteById(ctx context.Context, token string) error
		IsUserExists(ctx context.Context, emailOrUsername string) error
	}
)
