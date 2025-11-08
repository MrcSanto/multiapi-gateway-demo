package usecase

import (
	"errors"
	"fmt"
	"go-api/model"
	"go-api/repository"

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

func (uu *UserUseCase) GetUsers() ([]model.User, error) {
	return uu.repository.GetUsers()
}

func (uc *UserUseCase) DeleteUser(id int) error {
	rowsAffected, err := uc.repository.DeleteUser(id)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuario nao encontrado")
	}

	return nil
}

func (uc *UserUseCase) UpdateUser(id int, user model.User) (*model.User, error) {
	err := uc.repository.UpdateUser(id, user)
	if err != nil {
		return nil, err
	}

	updatedUser, err := uc.repository.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (uu *UserUseCase) GetUserById(id int) (*model.User, error) {
	user, err := uu.repository.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uu *UserUseCase) CreateUser(user model.User) (model.User, error) {
	// validando criação de usuario
	if user.Name == "" {
		return model.User{}, errors.New("nome é obrigatório")
	}
	if user.Email == "" {
		return model.User{}, errors.New("email é obrigatório")
	}
	if user.User == "" {
		return model.User{}, errors.New("username é obrigatório")
	}
	if user.Password == "" {
		return model.User{}, errors.New("senha é obrigatória")
	}

	// Verifica se o email já existe
	existingUser, err := uu.repository.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return model.User{}, errors.New("email já está em uso")
	}

	// Verifica se o username já existe
	existingUser, err = uu.repository.GetUserByUserName(user.User)
	if err == nil && existingUser != nil {
		return model.User{}, errors.New("username já está em uso")
	}

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
