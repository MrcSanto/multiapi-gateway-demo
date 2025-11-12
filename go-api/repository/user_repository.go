package repository

import (
	"go-api/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (ur *UserRepository) GetUsers() ([]model.User, error) {
	var users []model.User
	result := ur.db.Order("created_at DESC").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (ur *UserRepository) GetUserById(idUser int) (*model.User, error) {
	var user model.User
	result := ur.db.First(&user, idUser)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := ur.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

func (ur *UserRepository) GetUserByUserName(username string) (*model.User, error) {
	var user model.User
	result := ur.db.Where("username = ?", username).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

func (ur *UserRepository) DeleteUser(idUser int) (int64, error) {
	result := ur.db.Delete(&model.User{}, idUser)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (ur *UserRepository) UpdateUser(idUser int, user model.User) error {
	result := ur.db.Model(&model.User{}).
		Where("id = ?", idUser).
		Updates(map[string]interface{}{
			"name":     user.Name,
			"email":    user.Email,
			"username": user.User,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (ur *UserRepository) CreateUser(user model.User) (int, error) {
	result := ur.db.Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(user.ID), nil
}
