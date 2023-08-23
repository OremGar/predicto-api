package configuraciones

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/OremGar/predicto-api/modelos"
	"github.com/golang-jwt/jwt/v5"
)

var (
	llavePrivada *rsa.PrivateKey
	llavePublica *rsa.PublicKey
)

func init() {
	bytesPrivada, err := os.ReadFile("./llaves/private_key.pem")
	if err != nil {
		log.Fatal("No se pudo leer el archivo privado: ", err)
	}

	bytesPublica, err := os.ReadFile("./llaves/public_key.pem")
	if err != nil {
		log.Fatal("No se pudo leer el archivo público: ", err)
	}

	llavePrivada, err = jwt.ParseRSAPrivateKeyFromPEM(bytesPrivada)
	if err != nil {
		log.Fatal("No se pudo hacer el parse a privatekey: ", err)
	}

	llavePublica, err = jwt.ParseRSAPublicKeyFromPEM(bytesPublica)
	if err != nil {
		log.Fatal("No se pudo hacer el parse a publickey: ", err)
	}
}

// Función para generar Json Web Tokens
func GenerarJWT(idUsuario int) (string, error) {
	var claims *modelos.Claims = &modelos.Claims{
		IdUsuario: idUsuario,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(360 * time.Minute)),
		},
	}

	var token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(llavePrivada)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidarJWT(token string) (bool, error, modelos.Claims) {
	var claims *modelos.Claims = &modelos.Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return llavePublica, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, fmt.Errorf("error al validar JWT, firma inválida"), modelos.Claims{}
		}
	}
	if !tkn.Valid {
		return false, fmt.Errorf("JWT no válido"), modelos.Claims{}
	}

	return true, nil, *claims
}
