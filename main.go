package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"github.com/Gmax76/identity/controller"
	"github.com/Gmax76/identity/database"
	"github.com/Gmax76/identity/middleware"
)

func init() {
}

func main() {
	var (
		db                       = database.NewSqliteDatabase()
		authenticationMiddleware = middleware.NewAuthenticationMiddleware()
		userController           = controller.NewUserController(db, authenticationMiddleware)
	)
	router := gin.Default()

	router.GET("/users", userController.GetAll)
	router.POST("/users", userController.Create)
	router.POST("/login", authenticationMiddleware.IsAuthenticated, userController.Login)
	router.Run("localhost:8080")
}
