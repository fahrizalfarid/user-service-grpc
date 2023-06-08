package repository

import (
	"context"

	"github.com/fahrizalfarid/user-service-grpc/src/model"
	"gorm.io/gorm"
)

type validatorRepository struct {
	Db *gorm.DB
}

func NewValidatorRepo(db *gorm.DB) model.ValidatorRepository {
	return &validatorRepository{
		Db: db,
	}
}

func (v *validatorRepository) IsUsernameExists(ctx context.Context, username string) bool {
	var exists bool

	v.Db.WithContext(ctx).
		Raw(`SELECT EXISTS
			(
			SELECT
				1 
			FROM
				user_profiles up
				JOIN user_credentials uc ON up.ID = uc.ID 
			WHERE
				uc.username = $1 
			AND up.deleted_at = 0 
		);`, username).Scan(&exists)

	return exists
}

func (v *validatorRepository) IsEmailExists(ctx context.Context, email string) bool {
	var exists bool

	v.Db.WithContext(ctx).
		Raw(`SELECT EXISTS
				(
				SELECT
					1 
				FROM
					user_profiles up
					JOIN user_credentials uc ON up.ID = uc.ID 
				WHERE
					up.email = $1 
				AND up.deleted_at = 0 
				);`, email).Scan(&exists)

	return exists
}

func (v *validatorRepository) IsUserExists(ctx context.Context, usernameOrEmail string) bool {
	var exists bool

	v.Db.WithContext(ctx).
		Raw(`SELECT EXISTS
				(
				SELECT
					1 
				FROM
					user_profiles up
					JOIN user_credentials uc ON up.ID = uc.ID 
				WHERE
					(up.email = $1 OR uc.username = $1)
				AND up.deleted_at = 0 );`, usernameOrEmail).Scan(&exists)

	return exists
}
