package configuraciones

import (
	"crypto/rsa"
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
	bytesPrivada, err := os.ReadFile("./llaves/jwt_private.pem")
	if err != nil {
		log.Fatal("No se pudo leer el archivo privado: ", err)
	}

	bytesPublica, err := os.ReadFile("./llaves/jwt_public.pem")
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

// todo: --------------------------------------------
// todo: funcion para generar el token del usuario
// todo: --------------------------------------------
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
