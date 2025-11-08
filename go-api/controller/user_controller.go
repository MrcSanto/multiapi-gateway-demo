package controller

import (
	"go-api/model"
	"go-api/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userUseCase usecase.UserUseCase
}

func NewUserController(usecase usecase.UserUseCase) userController {
	return userController{
		userUseCase: usecase,
	}
}

func (u *userController) GetUsers(ctx *gin.Context) {
	users, err := u.userUseCase.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (u *userController) GetUserById(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id deve ser um inteiro",
		})
		return
	}

	user, err := u.userUseCase.GetUserById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "usuario nao encontrado",
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (u *userController) DeleteUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id precisa ser um inteiro"})
		return
	}

	err = u.userUseCase.DeleteUser(id)
	if err != nil {
		if err.Error() == "usuario nao encontrado" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "usuario nao encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao deletar o usuario"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "sucesso ao deletar o usuario"})
}

func (u *userController) UpdateUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id precisa ser um inteiro"})
		return
	}

	var user model.User
	err = ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	updatedUser, err := u.userUseCase.UpdateUser(id, user)
	if err != nil {
		if err.Error() == "usuario nao encontrado" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "usuario nao encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao atualizar o usuario"})
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}

func (u *userController) CreateUser(ctx *gin.Context) {
	var user model.User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	insertedUser, err := u.userUseCase.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	insertedUser.Password = ""

	ctx.JSON(http.StatusCreated, insertedUser)
}
