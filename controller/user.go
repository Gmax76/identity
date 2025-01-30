package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Gmax76/identity/database"
	"github.com/Gmax76/identity/entity"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserController interface {
		GetAll(ctx *gin.Context)
		Create(ctx *gin.Context)
		Login(ctx *gin.Context)
	}
	userController struct {
		db database.Database
	}
)

var secretKey = []byte("your-secret-key")

func NewUserController(db database.Database) UserController {
	return &userController{
		db: db,
	}
}

func (c *userController) GetAll(ctx *gin.Context) {
	users, err := c.db.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("error: %v", err))
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

func createToken(username string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"iss": "auth-app",                       // Issuer
		"aud": "user",                           // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (c *userController) Login(ctx *gin.Context) {
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

		tokenString, err := createToken(user.Email)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Error creating token")
			return
		}
		ctx.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
		ctx.JSON(http.StatusOK, tokenString)
		return
	}
}
