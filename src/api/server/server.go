package server

import (
	"github.com/fahrizalfarid/user-service-grpc/conf"
	v1 "github.com/fahrizalfarid/user-service-grpc/src/api/delivery/v1"
	"github.com/fahrizalfarid/user-service-grpc/src/api/dialer"
	"github.com/fahrizalfarid/user-service-grpc/src/api/middleware"
	"github.com/fahrizalfarid/user-service-grpc/src/api/usecase"
	"github.com/fahrizalfarid/user-service-grpc/utils"
	"github.com/labstack/echo/v4"

	mid "github.com/labstack/echo/v4/middleware"
)

func RunServer() *echo.Echo {
	conf.LoadEnv("./env")

	e := echo.New()
	e.Use(mid.Recover())

	auth := utils.NewAuthentication()
	dialer := dialer.NewGrpcDialer()

	userGrpcClientAddr := conf.GetUserClient()
	userValGrpcClientAddr := conf.GetValidatorClient()

	authorization := middleware.NewAuthorizationMiddleware(auth)
	userUsecase := usecase.NewUserUsecase(auth, dialer, userGrpcClientAddr, userValGrpcClientAddr)

	v1.NewUserDelivery(e, userUsecase, authorization)

	return e
}
