package main

import (
	"go-api/controller"
	"go-api/db"
	"go-api/repository"
	"go-api/usecase"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// carregando a .env (opcional)
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado, usando variáveis padrão")
	}

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

	// criando grupo para a api em go
	goGroup := server.Group("/go")

	// healthcheck endpoint
	goGroup.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": true,
		})
	})

	goGroup.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "API em Golang funcionando!")
	})

	goGroup.GET("/users", UserController.GetUsers)
	goGroup.GET("/users/:id", UserController.GetUserById)
	goGroup.POST("/users", UserController.CreateUser)
	goGroup.DELETE("/users/:id", UserController.DeleteUser)
	goGroup.PUT("/users/:id", UserController.UpdateUser)

	server.Run(":8000")
}
