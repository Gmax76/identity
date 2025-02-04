package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Gmax76/identity/database"
	"github.com/Gmax76/identity/entity"
	"github.com/Gmax76/identity/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserController interface {
		GetAll(ctx *gin.Context)
		Create(ctx *gin.Context)
		Login(ctx *gin.Context)
	}
	userController struct {
		db             database.Database
		authMiddleWare middleware.AuthenticationMiddleware
	}
)

func NewUserController(db database.Database, am middleware.AuthenticationMiddleware) UserController {
	return &userController{
		db:             db,
		authMiddleWare: am,
	}
}

func (c *userController) GetAll(ctx *gin.Context) {
	users, err := c.db.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("error: %v", err))
		return
	}
	ctx.IndentedJSON(http.StatusOK, users)
}

func (c *userController) Create(ctx *gin.Context) {
	var newUser entity.User
	if err := ctx.BindJSON(&newUser); err != nil {
		log.Printf("Error: %v", err)
		return
	}
	hashedPass, err := bcrypt.GenerateFromPassword(newUser.Password, bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("error: %v", err))
		return
	}
	newUser.Password = hashedPass
	createdUser, err := c.db.CreateUser(newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("error: %v", err))
		return
	}
	ctx.IndentedJSON(http.StatusCreated, createdUser)
}

func (c *userController) Login(ctx *gin.Context) {
	auth, authSet := ctx.Get("authenticated")
	if !authSet {
		log.Print("Pas de header authorization")
	} else {
		log.Printf("Auth: %v", auth)
	}
	var user entity.User
	var hashedPass []byte
	if err := ctx.BindJSON(&user); err != nil {
		log.Printf("Error: %v", err)
		return
	}
	dbUser, err := c.db.GetUser(user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, user.Email)
		return
	}
	hashedPass = dbUser.Password
	err = bcrypt.CompareHashAndPassword(hashedPass, user.Password)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		ctx.JSON(http.StatusUnauthorized, user.Email)
		return
	}
	if err == nil {

		tokenString, err := c.authMiddleWare.CreateToken(user.Email)
		if err != nil {
			log.Print(err)
			ctx.String(http.StatusInternalServerError, "Error creating token")
			return
		}
		ctx.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
		ctx.String(http.StatusOK, tokenString)
		return
	}
}
