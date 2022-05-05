package util

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("my_salt")

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	User
	jwt.StandardClaims
}

func GenerateToken(user User) (string, error) {
	expireTime := time.Now().Add(time.Minute * 30)

	claims := Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin_demo",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
