package main

import (
	_ "github.com/go-playground/validator/v10"
	_ "github.com/goccy/go-json"
	_ "github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/joho/godotenv"
	_ "github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4/middleware"
	_ "github.com/stretchr/testify"
	_ "github.com/stretchr/testify/suite"
	_ "go-micro.dev/v4"
	_ "go-micro.dev/v4/api"
	_ "golang.org/x/crypto/bcrypt"
	_ "google.golang.org/grpc"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
)
