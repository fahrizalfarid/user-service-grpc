package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/fahrizalfarid/user-service-grpc/conf"
	pb "github.com/fahrizalfarid/user-service-grpc/src/proto"
	"github.com/fahrizalfarid/user-service-grpc/src/user-service/delivery"
	"github.com/fahrizalfarid/user-service-grpc/src/user-service/repository"
	"github.com/fahrizalfarid/user-service-grpc/src/user-service/usecase"
	"github.com/fahrizalfarid/user-service-grpc/utils"
	"google.golang.org/grpc"
)

func RunUserSrv(port string) error {
	db, err := conf.DatabaseConn()
	if err != nil {
		panic(err)
	}

	auth := utils.NewAuthentication()

	userRepo := repository.NewUserRepo(db)
	userUsecase := usecase.NewUserSvcUsecase(userRepo)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	user := &delivery.User{
		UserUsecase:    userUsecase,
		Authentication: auth,
		Mu:             sync.Mutex{},
	}

	pb.RegisterUserServer(s, user)

	return s.Serve(lis)
}
