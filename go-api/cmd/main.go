package main

import (
	"go-api/controller"
	"go-api/db"
	"go-api/middleware"
	"go-api/repository"
	"go-api/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	// conexão com o banco de dados
	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	// camada de repo
	UserRepository := repository.NewUserRepository(dbConnection)
	ProductRepository := repository.NewProductRepository(dbConnection)

	// camada de usecase
	UserUseCase := usecase.NewUserUseCase(UserRepository)
	ProductUseCase := usecase.NewProductUseCase(ProductRepository)

	// camada de controller
	UserController := controller.NewUserController(UserUseCase)
	ProductController := controller.NewProductController(ProductUseCase)

	// healthcheck endpoint
	server.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": true,
		})
	})

	server.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "API funcionando!")
	})

	server.POST("/signin", UserController.CreateUser)
	server.POST("/login", UserController.Login)

	privateRoutes := server.Group("/")
	privateRoutes.Use(middleware.AuthMiddleware())
	{
		privateRoutes.GET("/products", ProductController.GetProducts)
		privateRoutes.GET("/products/:productId", ProductController.GetProductById)
		privateRoutes.POST("/product", ProductController.CreateProduct)
	}

	server.Run(":8000")
}
