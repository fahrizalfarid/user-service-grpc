package delivery

import (
	"context"

	"github.com/fahrizalfarid/user-service-grpc/src/model"
	pb "github.com/fahrizalfarid/user-service-grpc/src/proto"
)

type UserValidator struct {
	pb.UnimplementedUserValidatorServer
	UserValidator model.UserValidatorSvcUsecase
}

func (u *UserValidator) IsUsernameExists(ctx context.Context, req *pb.UsernameRequest) (*pb.Found, error) {
	found := new(pb.Found)
	exists := u.UserValidator.IsUsernameExists(ctx, req.Username)
	found.Found = exists
	return found, nil
}

func (u *UserValidator) IsEmailExists(ctx context.Context, req *pb.EmailRequest) (*pb.Found, error) {
	found := new(pb.Found)
	exists := u.UserValidator.IsEmailExists(ctx, req.Email)
	found.Found = exists
	return found, nil
}

func (u *UserValidator) IsUserExists(ctx context.Context, req *pb.EmailOrUsernameRequest) (*pb.Found, error) {
	found := new(pb.Found)
	exists := u.UserValidator.IsUserExists(ctx, req.EmailOrUsername)
	found.Found = exists
	return found, nil
}
