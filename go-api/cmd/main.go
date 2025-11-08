package main

import (
	"go-api/controller"
	"go-api/db"
	"go-api/repository"
	"go-api/usecase"
	"os"

	"github.com/gin-gonic/gin"
)

var INSTANCE_ID = getInstanceID()

func getInstanceID() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	return hostname
}

func main() {
	server := gin.Default()

	// conexão com o banco de dados
	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	// camada de repo
	UserRepository := repository.NewUserRepository(dbConnection)

	// camada de usecase
	UserUseCase := usecase.NewUserUseCase(UserRepository)

	// camada de controller
	UserController := controller.NewUserController(UserUseCase)

	// healthcheck endpoint
	server.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": true,
		})
	})

	server.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "API funcionando!")
	})

	server.GET("/users", UserController.GetUsers)
	server.GET("/users/:id", UserController.GetUserById)
	server.POST("/users", UserController.CreateUser)
	server.DELETE("/users/:id", UserController.DeleteUser)
	server.PUT("/users/:id", UserController.UpdateUser)

	server.Run(":8000")
}
