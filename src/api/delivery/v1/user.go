package v1

import (
	"net/http"
	"strings"

	"github.com/fahrizalfarid/user-service-grpc/src/api/middleware"
	"github.com/fahrizalfarid/user-service-grpc/src/model"
	"github.com/fahrizalfarid/user-service-grpc/src/request"
	"github.com/fahrizalfarid/user-service-grpc/src/response"
	"github.com/fahrizalfarid/user-service-grpc/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type userDelivery struct {
	UserUsecase model.UserUsecase
	Middleware  middleware.Authorization
	Validator   *validator.Validate
}

func NewUserDelivery(
	e *echo.Echo,
	userUsecase model.UserUsecase,
	auth middleware.Authorization,
) {
	handler := &userDelivery{
		UserUsecase: userUsecase,
		Middleware:  auth,
		Validator:   validator.New(),
	}

	e.Validator = handler
	v1 := e.Group("/api/v1")

	v1.POST("/user/create", handler.Create)
	v1.GET("/user", handler.GetById, handler.Middleware.ValidateToken)
	v1.POST("/user/login", handler.Login)
	v1.POST("/user/update", handler.UpdateById, handler.Middleware.ValidateToken)
	v1.GET("/user/find", handler.Find, handler.Middleware.ValidateToken)
	v1.POST("/user/delete", handler.DeleteUser, handler.Middleware.ValidateToken)
}

func (u *userDelivery) Validate(v any) error {
	return u.Validator.Struct(v)
}

func (u *userDelivery) Create(e echo.Context) error {
	var userRequest *request.UserRequest

	if err := e.Bind(&userRequest); err != nil {
		return e.JSON(http.StatusBadRequest, utils.BadRequest(err.Error()))
	}

	if err := u.Validate(userRequest); err != nil {
		return e.JSON(http.StatusBadRequest, utils.BadRequest(err.Error()))
	}

	ctx := e.Request().Context()
	token, id, err := u.UserUsecase.CreateUser(ctx, userRequest)
	if err != nil {
		return e.JSON(http.StatusBadRequest, utils.BadRequest(err.Error()))
	}

	return e.JSON(http.StatusOK, utils.SuccessResponse(response.AuthResponse{
		Token:    token,
		Username: userRequest.Username,
		UserId:   id,
	}))
}

func (u *userDelivery) GetById(e echo.Context) error {
	id := u.Middleware.GetTokenData().Id

	ctx := e.Request().Context()
	user, err := u.UserUsecase.GetUserById(ctx, id, u.Middleware.GetToken())

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return e.JSON(http.StatusNotFound, utils.FailResponseWithId(id))
		}
		return e.JSON(http.StatusBadRequest, utils.FailResponse(err.Error()))
	}

	return e.JSON(http.StatusOK, utils.SuccessResponse(user))
}

func (u *userDelivery) Login(e echo.Context) error {
	var request *request.UserLoginRequest

	if err := e.Bind(&request); err != nil {
		return e.JSON(http.StatusBadRequest, utils.BadRequest(err.Error()))
	}

	if err := u.Validate(request); err != nil {
		return e.JSON(http.StatusBadRequest, utils.BadRequest(err.Error()))
	}

	ctx := e.Request().Context()
	authResponse, err := u.UserUsecase.Login(ctx, request.Username, request.Password)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return e.JSON(http.StatusNotFound, utils.FailResponseWithId(request.Username))
		}
		return e.JSON(http.StatusUnauthorized, utils.FailResponse(err.Error()))
	}

	return e.JSON(http.StatusOK, utils.SuccessResponse(&authResponse))
}

func (u *userDelivery) Find(e echo.Context) error {
	id := e.QueryParam("query")
	if id == "" {
		return e.JSON(http.StatusBadRequest, utils.BadRequest("missing username or email"))
	}

	ctx := e.Request().Context()
	// stream
	userProfile, err := u.UserUsecase.Find(ctx, id, u.Middleware.GetToken())

	// userProfile, err := u.UserUsecase.FindWithArray(ctx, id, u.Middleware.GetToken())
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return e.JSON(http.StatusNotFound, utils.FailResponseWithId(id))
		}
		return e.JSON(http.StatusBadRequest, utils.FailResponse(err.Error()))
	}

	return e.JSON(http.StatusOK, utils.SuccessResponse(userProfile))
}

func (u *userDelivery) UpdateById(e echo.Context) error {
	var request *request.UserUpdateRequest

	if err := e.Bind(&request); err != nil {
		return e.JSON(http.StatusBadRequest, utils.BadRequest(err.Error()))
	}

	id := u.Middleware.GetTokenData().Id
	if err := u.Validate(request); err != nil {
		return e.JSON(http.StatusBadRequest, utils.BadRequest(err.Error()))
	}

	ctx := e.Request().Context()
	userProfile, err := u.UserUsecase.UpdateById(ctx, request, u.Middleware.GetToken())
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return e.JSON(http.StatusNotFound, utils.FailResponseWithId(id))
		}
		return e.JSON(http.StatusBadRequest, utils.FailResponse(err.Error()))
	}

	return e.JSON(http.StatusOK, utils.SuccessResponse(userProfile))
}

func (u *userDelivery) DeleteUser(e echo.Context) error {
	userId := u.Middleware.GetTokenData().Id

	ctx := e.Request().Context()
	err := u.UserUsecase.DeleteById(ctx, u.Middleware.GetToken())

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return e.JSON(http.StatusNotFound, utils.FailResponseWithId(userId))
		}
		return e.JSON(http.StatusBadRequest, utils.BadRequest(err.Error()))
	}

	return e.JSON(http.StatusOK, utils.SuccesNilsResponse())
}
