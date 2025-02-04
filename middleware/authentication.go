package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthenticationMiddleware interface {
	IsAuthenticated(ctx *gin.Context)
	CreateToken(username string) (string, error)
}

type authenticationMiddleware struct {
	secretKey        []byte
	jwtSigningMethod jwt.SigningMethod
}

func NewAuthenticationMiddleware() AuthenticationMiddleware {
	return &authenticationMiddleware{
		secretKey:        []byte("azertyuiop"),
		jwtSigningMethod: jwt.SigningMethodHS256,
	}
}

func (m *authenticationMiddleware) CreateToken(username string) (string, error) {
	claims := jwt.NewWithClaims(m.jwtSigningMethod, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"iss": "auth-app",                       // Issuer
		"aud": "user",                           // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(m.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *authenticationMiddleware) IsAuthenticated(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	token := strings.Split(authHeader, " ")
	if len(token) > 1 {
		token, err := jwt.Parse(token[1], func(token *jwt.Token) (interface{}, error) {
			return m.secretKey, nil
		})
		if err != nil {
			log.Print(err)
			log.Print("Error parsing token")
			ctx.Next()
			return
		}
		if !token.Valid {
			log.Print("Token invalid")
			ctx.Next()
			return
		}
		ctx.Set("authenticated", true)
		return
	}
	ctx.Next()
}
