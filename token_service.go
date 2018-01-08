package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	pb "github.com/jakebjorke/shipper-user-service/proto/user"
)

var (
	key = []byte("this is just for testing so who really cares")
)

//CustomClaims is used to add application claims to a JWT
type CustomClaims struct {
	User *pb.User
	jwt.StandardClaims
}

//Authable is used for tokens
type Authable interface {
	Decode(tokenString string) (*CustomClaims, error)
	Encode(user *pb.User) (string, error)
}

//TokenService is used to implement the Authable interface
type TokenService struct {
	repo Repository
}

//Decode from string to an object
func (srv *TokenService) Decode(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

//Encode a claim into token
func (srv *TokenService) Encode(user *pb.User) (string, error) {
	expireToken := time.Now().Add(time.Hour * 72).Unix()

	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "go.micro.srv.user",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)

	return token.SignedString(key)
}
