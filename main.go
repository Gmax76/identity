package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	"github.com/Gmax76/identity/database"
)

var (
	secretKey = []byte("your-secret-key")
	db        database.Database
)

func getUsers(c *gin.Context) {
	users, err := db.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("error: %v", err))
	}
	c.IndentedJSON(http.StatusOK, users)
}

func postUsers(c *gin.Context) {
	var newUser database.User
	if err := c.BindJSON(&newUser); err != nil {
		log.Printf("Error: %v", err)
		return
	}
	hashedPass, err := bcrypt.GenerateFromPassword(newUser.Password, bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("error: %v", err))
		return
	}
	newUser.Password = hashedPass
	createdUser, err := db.CreateUser(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("error: %v", err))
		return
	}
	c.IndentedJSON(http.StatusCreated, createdUser)
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

func postLogin(c *gin.Context) {
	var user database.User
	var hashedPass []byte
	if err := c.BindJSON(&user); err != nil {
		log.Printf("Error: %v", err)
		return
	}
	dbUser, err := db.GetUser(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, user.Email)
		return
	}
	hashedPass = dbUser.Password
	err = bcrypt.CompareHashAndPassword(hashedPass, user.Password)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusUnauthorized, user.Email)
		return
	}
	if err == nil {

		tokenString, err := createToken(user.Email)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error creating token")
			return
		}
		c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
		c.JSON(http.StatusOK, tokenString)
		return
	}
}

func init() {
}

func main() {
	db = database.NewSqliteDatabase()

	router := gin.Default()

	router.GET("/users", getUsers)
	router.POST("/users", postUsers)
	router.POST("/login", postLogin)
	router.Run("localhost:8080")
}
