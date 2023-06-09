package delivery

import (
	"context"
	"sync"

	"github.com/fahrizalfarid/user-service-grpc/src/constant"
	"github.com/fahrizalfarid/user-service-grpc/src/model"
	pb "github.com/fahrizalfarid/user-service-grpc/src/proto"
	"github.com/fahrizalfarid/user-service-grpc/src/request"
	"github.com/fahrizalfarid/user-service-grpc/utils"
	"google.golang.org/grpc/metadata"
)

type User struct {
	pb.UnimplementedUserServer
	UserUsecase    model.UserSvcUsecase
	Authentication utils.Authentication
	Mu             sync.Mutex
}

func (u *User) getToken(ctx context.Context) (*utils.Token, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, constant.ErrAuth
	}

	token, exist := md["token"]
	if !exist {
		return nil, constant.ErrAuth
	}
	claims, err := u.Authentication.ParsingToken(token[0])
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (u *User) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	hashedPassword, err := u.Authentication.EncryptPassword(req.Password)
	if err != nil {
		return nil, err
	}

	id, err := u.UserUsecase.CreateUser(ctx, &request.UserRequestService{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Email:     req.Email,
		Address:   req.Address,
		Phone:     req.Phone,
		Username:  req.Username,
		Password:  hashedPassword,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{Id: id}, nil
}

func (u *User) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.UserResponse, error) {
	_, err := u.getToken(ctx)
	if err != nil {
		return nil, err
	}

	user, err := u.UserUsecase.GetUserById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
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

func (u *User) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	userCredential, err := u.UserUsecase.GetUserCredentials(ctx, req.UsernameOrEmail)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Id:       userCredential.Id,
		Username: userCredential.Username,
		Password: userCredential.Password,
	}, nil
}

func (u *User) Find(req *pb.FindRequest, stream pb.User_FindServer) error {

	_, err := u.getToken(stream.Context())
	if err != nil {
		return err
	}

	u.Mu.Lock()
	user, err := u.UserUsecase.GetUserByEmailOrUsername(stream.Context(), req.Word)
	if err != nil {
		return err
	}

	for _, v := range user {
		rsp := &pb.UserFound{
			Id:       v.Id,
			Username: v.Username,
			Fullname: v.Fullname,
			Email:    v.Email,
		}

		if err := stream.Send(rsp); err != nil {
			return err
		}
	}
	u.Mu.Unlock()

	return nil
}

func (u *User) FindWithArray(ctx context.Context, req *pb.FindRequest) (*pb.UserFoundArray, error) {
	_, err := u.getToken(ctx)
	if err != nil {
		return nil, err
	}

	var data []*pb.UserFound

	user, err := u.UserUsecase.GetUserByEmailOrUsername(ctx, req.Word)
	if err != nil {
		return nil, err
	}

	for _, v := range user {
		data = append(data, &pb.UserFound{
			Id:       v.Id,
			Username: v.Username,
			Fullname: v.Fullname,
			Email:    v.Email,
		})
	}

	return &pb.UserFoundArray{Users: data}, nil
}

func (u *User) UpdateById(ctx context.Context, req *pb.UpdateRequest) (*pb.UserResponse, error) {
	claims, err := u.getToken(ctx)
	if err != nil {
		return nil, err
	}

	user, err := u.UserUsecase.UpdateById(ctx, &request.UserUpdateService{
		Id:        claims.Id,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		Username:  req.Address,
		Password:  req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
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

func (u *User) DeleteById(ctx context.Context, req *pb.DeleteRequest) (*pb.Error, error) {
	claims, err := u.getToken(ctx)
	if err != nil {
		return nil, err
	}

	err = u.UserUsecase.DeleteById(ctx, claims.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Error{
		Message: "",
	}, nil
}
