package usecase

import (
	"errors"
	"go-api/model"
	"go-api/repository"
	"go-api/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repository repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return UserUseCase{
		repository: repo,
	}
}

func (uu *UserUseCase) LoginUser(email, password string) (string, error) {
	user, err := uu.repository.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("usuário não encontrado")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("senha incorreta")
	}

	tokenString, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (uu *UserUseCase) CreateUser(user model.User) (model.User, error) {

	// fazendo o hash da senha
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	user.Password = string(hashPassword)

	userId, err := uu.repository.CreateUser(user)
	if err != nil {
		return model.User{}, err
	}

	user.ID = userId

	return user, nil
}
