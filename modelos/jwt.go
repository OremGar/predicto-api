package modelos

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	IdUsuario int
	jwt.RegisteredClaims
}
