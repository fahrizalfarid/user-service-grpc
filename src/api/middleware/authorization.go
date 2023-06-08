package middleware

import (
	"net/http"
	"strings"

	"github.com/fahrizalfarid/user-service-grpc/utils"
	"github.com/labstack/echo/v4"
)

type authorization struct {
	Authentication utils.Authentication
	TokenData      *utils.Token
	Token          *string
}

type Authorization interface {
	ValidateToken(next echo.HandlerFunc) echo.HandlerFunc
	GetTokenData() *utils.Token
	GetToken() string
}

func NewAuthorizationMiddleware(auth utils.Authentication) Authorization {
	return &authorization{
		Authentication: auth,
	}
}

func (a *authorization) ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		header := e.Request().Header
		authorization := header.Get("Authorization")

		if !strings.Contains(authorization, "Bearer") {
			return e.JSON(http.StatusBadRequest, utils.BadRequest("Authorization header is required"))
		}

		token := strings.Replace(authorization, "Bearer ", "", -1)
		if token == "" || token == "Bearer" {
			return e.JSON(http.StatusBadRequest, utils.BadRequest("Bearer token is required"))
		}

		tokenData, err := a.Authentication.ParsingToken(token)
		if err != nil {
			return e.JSON(http.StatusUnauthorized, utils.BadRequest(err.Error()))
		}

		a.TokenData = tokenData
		a.Token = &token
		return next(e)
	}
}

func (a *authorization) GetTokenData() *utils.Token {
	return a.TokenData
}

func (a *authorization) GetToken() string {
	return *a.Token
}
