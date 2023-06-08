package server

import (
	"fmt"
	"net"

	"github.com/fahrizalfarid/user-service-grpc/conf"
	pb "github.com/fahrizalfarid/user-service-grpc/src/proto"
	"github.com/fahrizalfarid/user-service-grpc/src/validator-service/delivery"
	"github.com/fahrizalfarid/user-service-grpc/src/validator-service/repository"
	"github.com/fahrizalfarid/user-service-grpc/src/validator-service/usecase"
	"google.golang.org/grpc"
)

func RunUserValidatorSrv(port string) error {

	db, err := conf.DatabaseConn()
	if err != nil {
		panic(err)
	}

	validator := repository.NewValidatorRepo(db)
	userValidatorUsecase := usecase.NewUserValidatorSvcUsecase(validator)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	val := &delivery.UserValidator{
		UserValidator: userValidatorUsecase,
	}

	pb.RegisterUserValidatorServer(s, val)
	return s.Serve(lis)
}
