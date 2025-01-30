package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"github.com/Gmax76/identity/controller"
	"github.com/Gmax76/identity/database"
)

func init() {
}

func main() {
	var (
		db             database.Database         = database.NewSqliteDatabase()
		userController controller.UserController = controller.NewUserController(db)
	)
	router := gin.Default()

	router.GET("/users", userController.GetAll)
	router.POST("/users", userController.Create)
	router.POST("/login", userController.Login)
	router.Run("localhost:8080")
}
