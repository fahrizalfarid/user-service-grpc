package usecase

import (
	"context"
	"io"

	"strings"

	"time"

	"github.com/fahrizalfarid/user-service-grpc/src/api/dialer"
	"github.com/fahrizalfarid/user-service-grpc/src/constant"
	"github.com/fahrizalfarid/user-service-grpc/src/model"
	pb "github.com/fahrizalfarid/user-service-grpc/src/proto"
	"github.com/fahrizalfarid/user-service-grpc/src/request"
	"github.com/fahrizalfarid/user-service-grpc/src/response"
	"github.com/fahrizalfarid/user-service-grpc/utils"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type userUsecase struct {
	GrpcDialer              dialer.Dialer
	UserSrvAddress          string
	UserValidatorSrvAddress string
	Authentication          utils.Authentication
}

func NewUserUsecase(auth utils.Authentication, d dialer.Dialer, userAddr, valAddr string) model.UserUsecase {
	return &userUsecase{
		GrpcDialer:              d,
		UserSrvAddress:          userAddr,
		UserValidatorSrvAddress: valAddr,
		Authentication:          auth,
	}
}

func (u *userUsecase) existsValidator(ctx context.Context, username, email string) error {
	conn, val, err := u.GrpcDialer.UserValidatorClient(u.UserValidatorSrvAddress)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	userFound, err := val.IsUsernameExists(ctx, &pb.UsernameRequest{Username: username})

	if err != nil {
		if strings.Contains(err.Error(), "Internal Server Error") {
			return constant.ErrSrvNotAvailable
		}
		return err
	}

	if userFound.Found {
		return constant.ErrUsernameExists
	}

	ctx, cancel = context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	emailFound, err := val.IsEmailExists(ctx, &pb.EmailRequest{
		Email: email,
	})

	if err != nil {
		if strings.Contains(err.Error(), "Internal Server Error") {
			return constant.ErrSrvNotAvailable
		}
		return err
	}

	if emailFound.Found {
		return constant.ErrEmailExists
	}
	return nil
}

func (u *userUsecase) CreateUser(ctx context.Context, data *request.UserRequest) (string, int64, error) {
	err := u.existsValidator(ctx, data.Username, data.Email)
	if err != nil {
		return "", 0, err
	}

	conn, userSrv, err := u.GrpcDialer.UserClient(u.UserSrvAddress)
	if err != nil {
		return "", 0, err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	rsp, err := userSrv.Create(ctx, &pb.CreateRequest{
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Email:     data.Email,
		CreatedAt: time.Now().Unix(),
		Phone:     data.Phone,
		Address:   data.Address,
		DeletedAt: int64(0),
		Username:  data.Username,
		Password:  data.Password,
	})

	if err != nil {
		if strings.Contains(err.Error(), "Internal Server Error") {
			return "", 0, constant.ErrSrvNotAvailable
		}
		return "", 0, err
	}

	stringToken, err := u.Authentication.GenerateToken(rsp.Id, data.Username)
	if err != nil {
		return "", 0, err
	}

	return stringToken, rsp.Id, nil
}

func (u *userUsecase) GetUserById(ctx context.Context, id int64, token string) (*response.UserProfileResponse, error) {
	md := metadata.Pairs("token", token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	conn, userSrv, err := u.GrpcDialer.UserClient(u.UserSrvAddress)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	user, err := userSrv.GetById(ctx, &pb.GetByIdRequest{Id: id})

	if err != nil {
		if strings.Contains(err.Error(), "Internal Server Error") {
			return nil, constant.ErrSrvNotAvailable
		}
		return nil, err
	}

	return &response.UserProfileResponse{
		Id:        user.Id,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Phone:     user.Phone,
		Address:   user.Address,
		Username:  user.Username,
	}, nil
}

func (u *userUsecase) Login(ctx context.Context, username, password string) (*response.AuthResponse, error) {
	conn, userSrv, err := u.GrpcDialer.UserClient(u.UserSrvAddress)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	user, err := userSrv.Login(ctx, &pb.LoginRequest{
		UsernameOrEmail: username,
	})

	if err != nil {
		if strings.Contains(err.Error(), "Internal Server Error") {
			return nil, constant.ErrSrvNotAvailable
		}
		return nil, err
	}

	err = u.Authentication.CompareHashAndPassword(user.Password, password)
	if err != nil {
		return nil, err
	}

	token, err := u.Authentication.GenerateToken(user.Id, user.Username)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		UserId:   user.Id,
		Username: user.Username,
		Token:    token,
	}, nil
}

func (u *userUsecase) Find(ctx context.Context, username, token string) ([]*response.UserFound, error) {
	var data []*response.UserFound

	md := metadata.Pairs("token", token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	conn, userSrv, err := u.GrpcDialer.UserClient(u.UserSrvAddress)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	connVall, val, err := u.GrpcDialer.UserValidatorClient(u.UserValidatorSrvAddress)
	if err != nil {
		return nil, err
	}
	defer connVall.Close()

	claims, _ := u.Authentication.ParsingToken(token)

	ctxT, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	found, err := val.IsUserExists(ctx, &pb.EmailOrUsernameRequest{EmailOrUsername: claims.Username})
	if err != nil {
		return nil, err
	}

	if !found.Found {
		return nil, constant.ErrUsernameNotExists
	}

	ctxT, cancel = context.WithTimeout(ctxT, 2*time.Second)
	defer cancel()

	user, err := userSrv.Find(ctxT, &pb.FindRequest{Word: username})

	if err != nil {
		if strings.Contains(err.Error(), "Internal Server Error") {
			return nil, constant.ErrSrvNotAvailable
		}
		return nil, err
	}

	for {
		rsp, err := user.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		data = append(data, &response.UserFound{
			Id:       rsp.Id,
			Username: rsp.Username,
			Fullname: rsp.Fullname,
			Email:    rsp.Email,
		})
	}

	if len(data) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return data, nil
}

func (u *userUsecase) FindWithArray(ctx context.Context, username, token string) (*pb.UserFoundArray, error) {
	md := metadata.Pairs("token", token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	conn, userSrv, err := u.GrpcDialer.UserClient(u.UserSrvAddress)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	connVall, val, err := u.GrpcDialer.UserValidatorClient(u.UserValidatorSrvAddress)
	if err != nil {
		return nil, err
	}
	defer connVall.Close()

	claims, _ := u.Authentication.ParsingToken(token)

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	found, err := val.IsUserExists(ctx, &pb.EmailOrUsernameRequest{EmailOrUsername: claims.Username})
	if err != nil {
		return nil, err
	}

	if !found.Found {
		return nil, constant.ErrUsernameNotExists
	}

	ctxT, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	user, err := userSrv.FindWithArray(ctxT, &pb.FindRequest{Word: username})

	if err != nil {
		if strings.Contains(err.Error(), "Internal Server Error") {
			return nil, constant.ErrSrvNotAvailable
		}
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) UpdateById(ctx context.Context, data *request.UserUpdateRequest, token string) (*response.UserProfileResponse, error) {
	md := metadata.Pairs("token", token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	conn, userSrv, err := u.GrpcDialer.UserClient(u.UserSrvAddress)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if data.Password != "" {
		hashedPassword, errHash := u.Authentication.EncryptPassword(data.Password)
		if errHash != nil {
			return nil, errHash
		}

		data.Password = hashedPassword
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := userSrv.UpdateById(ctx, &pb.UpdateRequest{
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Email:     data.Email,
		Phone:     data.Phone,
		Address:   data.Address,
		Username:  data.Username,
		Password:  data.Password,
	})

	if err != nil {
		if strings.Contains(err.Error(), "Internal Server Error") {
			return nil, constant.ErrSrvNotAvailable
		}
		return nil, err
	}

	return &response.UserProfileResponse{
		Id:        user.Id,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Phone:     user.Phone,
		Address:   user.Address,
		Username:  user.Username,
	}, nil
}

func (u *userUsecase) DeleteById(ctx context.Context, token string) error {
	md := metadata.Pairs("token", token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	conn, userSrv, err := u.GrpcDialer.UserClient(u.UserSrvAddress)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctxT, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err = userSrv.DeleteById(ctxT, &pb.DeleteRequest{})
	if err != nil {
		if strings.Contains(err.Error(), "Internal Server Error") {
			return constant.ErrSrvNotAvailable
		}
		return err
	}
	return nil
}

func (u *userUsecase) IsUserExists(ctx context.Context, emailOrUsername string) error {
	conn, val, err := u.GrpcDialer.UserValidatorClient(u.UserValidatorSrvAddress)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	found, err := val.IsUserExists(ctx, &pb.EmailOrUsernameRequest{EmailOrUsername: emailOrUsername})
	if err != nil {
		if strings.Contains(err.Error(), "Internal Server Error") {
			return constant.ErrSrvNotAvailable
		}
		return err
	}

	if !found.Found {
		return gorm.ErrRecordNotFound
	}

	return nil
}
