package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/fahrizalfarid/user-service-grpc/src/model"
	"github.com/fahrizalfarid/user-service-grpc/src/response"
	"gorm.io/gorm"
)

type userRepository struct {
	Db *gorm.DB
}

func NewUserRepo(db *gorm.DB) model.UserRepository {
	return &userRepository{
		Db: db,
	}
}

func (u *userRepository) CreateUserProfile(ctx context.Context, data *model.UserProfiles) (int64, int64, error) {
	var result *model.UserProfiles

	err := u.Db.WithContext(ctx).Create(data).Scan(&result)

	if err.Error != nil || err.RowsAffected == 0 {
		return 0, 0, err.Error
	}

	return err.RowsAffected, result.Id, nil
}

func (u *userRepository) CreateUserCredential(ctx context.Context, dataCreds *model.UserCredentials) (int64, error) {
	err := u.Db.WithContext(ctx).Create(dataCreds)

	if err.Error != nil || err.RowsAffected == 0 {
		return 0, err.Error
	}

	return err.RowsAffected, nil
}

func (u *userRepository) GetUserById(ctx context.Context, id int64) (*response.UserProfile, error) {
	var data *response.UserProfile

	err := u.Db.WithContext(ctx).Raw(
		`
		SELECT
			up.id,
			up.firstname,
			up.lastname,
			up.email,
			up.created_at,
			up.phone,
			up.address,
			uc.username
		FROM
			user_profiles up
			JOIN user_credentials uc ON up.id = uc.id
		WHERE
		    up.deleted_at = 0 AND up.id = $1;
		`, id).Scan(&data)

	if err.Error != nil || err.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return data, nil
}

func (u *userRepository) GetUserCredentials(ctx context.Context, usernameOrEmail string) (*response.UserLoginResponse, error) {
	var result *response.UserLoginResponse

	err := u.Db.WithContext(ctx).Raw(
		`
		SELECT
			up.id,
			uc.username,
			uc."password"
		FROM
			user_profiles up
			JOIN user_credentials uc ON up."id" = uc."id"
		WHERE
			up.deleted_at = 0
			AND ((up.email = $1) OR (uc.username = $1));
		`, usernameOrEmail,
	).Scan(&result)

	if err.Error != nil || err.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return result, nil
}

func (u *userRepository) GetUserByUsernameOrEmail(ctx context.Context, id string) ([]*response.UserFound, error) {
	var data []*response.UserFound

	ilikeId := fmt.Sprintf("%%%s%%", id)

	err := u.Db.WithContext(ctx).Raw(
		`
		SELECT
			up.id,
			uc.username,
			concat (up.firstname, ' ', up.lastname) fullname,
			up.email
		FROM
			user_profiles up
			JOIN user_credentials uc ON up."id" = uc."id"
		WHERE
			up.deleted_at = 0
			AND ((up.email ILIKE $1) OR (uc.username ILIKE $1));
		`, ilikeId).Scan(&data)

	if err.Error != nil || err.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return data, nil
}

func (u *userRepository) UpdateUserById(ctx context.Context, data *model.UserProfiles) (int64, error) {
	err := u.Db.WithContext(ctx).
		Model(&model.UserProfiles{}).
		Where("id = ?", data.Id).
		Where("deleted_at = ?", int64(0)).
		Updates(data)

	if err.Error != nil || err.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return err.RowsAffected, nil
}

func (u *userRepository) UpdateUserCredentialById(ctx context.Context, data *model.UserCredentials) (int64, error) {
	err := u.Db.WithContext(ctx).Model(&model.UserCredentials{
		Id: data.Id,
	}).Updates(data)

	if err.Error != nil || err.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return err.RowsAffected, nil
}

func (u *userRepository) DeleteUserById(ctx context.Context, id int64) (int64, error) {
	err := u.Db.WithContext(ctx).
		Model(&model.UserProfiles{}).
		Where("deleted_at = ?", int64(0)).
		Where("id = ?", id).
		Update("deleted_at", time.Now().Unix())

	if err.Error != nil || err.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return err.RowsAffected, nil
}
