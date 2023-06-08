package conf

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadEnv(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}
	return nil
}

func GetUserSrv() string {
	return os.Getenv("USER_SRV")
}

func GetValidatorSrv() string {
	return os.Getenv("VALIDATOR_SRV")
}

func GetUserClient() string {
	return os.Getenv("USER_SRV_CLIENT")
}

func GetValidatorClient() string {
	return os.Getenv("VALIDATOR_SRV_CLIENT")
}

func DatabaseConn() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("HOST_POSTGRES"),
		os.Getenv("USERNAME_POSTGRES"),
		os.Getenv("PASSWORD_POSTGRES"),
		os.Getenv("DATABASE"),
		os.Getenv("PORT_POSTGRES"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()

	err = sqlDb.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	sqlDb.SetMaxIdleConns(0)
	sqlDb.SetMaxOpenConns(10)
	sqlDb.SetConnMaxLifetime(1 * time.Minute)

	return db, nil
}
