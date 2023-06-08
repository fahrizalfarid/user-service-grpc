package utils

import (
	"errors"
	"time"

	"github.com/fahrizalfarid/user-service-grpc/src/constant"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type authentication struct {
	SecretKey []byte
}

type Authentication interface {
	GenerateToken(userId int64, username string) (string, error)
	ParsingToken(tokenString string) (*Token, error)
	RefreshToken(tokenString string) (string, error)
	EncryptPassword(password string) (string, error)
	CompareHashAndPassword(hashedPassword, password string) error
}

func NewAuthentication() Authentication {
	return &authentication{
		SecretKey: []byte(constant.SigningKey),
	}
}

func (a *authentication) GenerateToken(userId int64, username string) (string, error) {
	claims := Token{
		Id:       userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(24 * 3 * time.Hour)},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(a.SecretKey)

	if err != nil {
		return "", err
	}

	return ss, err
}

func (a *authentication) ParsingToken(tokenString string) (*Token, error) {
	myClaims := &Token{}

	token, err := jwt.ParseWithClaims(tokenString, myClaims, func(t *jwt.Token) (interface{}, error) {
		return a.SecretKey, nil
	})

	if token.Valid {
		return myClaims, nil
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return nil, err
	} else if errors.Is(err, jwt.ErrSignatureInvalid) {
		return nil, err
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return nil, err
	} else {
		return nil, err
	}
}

func (a *authentication) RefreshToken(tokenString string) (string, error) {
	myClaims := &Token{}

	_, _ = jwt.ParseWithClaims(tokenString, myClaims, func(t *jwt.Token) (interface{}, error) {
		return a.SecretKey, nil
	})

	newExpired := &jwt.NumericDate{Time: time.Now().Add(24 * 3 * time.Hour)}

	myClaims.ExpiresAt = newExpired

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)
	ss, err := newToken.SignedString(a.SecretKey)

	if err != nil {
		return "", err
	}

	return ss, err
}

func (a *authentication) EncryptPassword(password string) (string, error) {
	bytePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}

func (a *authentication) CompareHashAndPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("wrong password")
	}

	return nil
}
