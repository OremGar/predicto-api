package modelos

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/OremGar/predicto-api/bd"
	"github.com/golang-jwt/jwt/v5"
)

type Claim struct {
	Usuario Usuarios `json:"usuario"`
	*jwt.RegisteredClaims
}

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func init() {
	privateBytes, err := os.ReadFile(".llaves/jwt_private.pem")
	if err != nil {
		log.Fatal("No se pudo leer el archivo privado: ", err)
	}

	publicBytes, err := os.ReadFile(".llaves/jwt_public.pem")
	if err != nil {
		log.Fatal("No se pudo leer el archivo público", err)
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal("No se pudo hacer el parse a privatekey")
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Fatal("No se pudo hacer el parse a publickey")
	}
}

// todo: --------------------------------------------
// todo: funcion para generar el token del usuario
// todo: --------------------------------------------
func GenerateJWT(usuario Usuarios) string {
	claims := &Claim{
		Usuario: usuario,
		RegisteredClaims: &jwt.RegisteredClaims{
			//240
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 360)),
			Issuer:    "Auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		log.Fatal("No se pudo firmar el token")
	}
	return tokenString
}

// Función que valida el usuario contenido en el token
func validateJWT(bearToken string) (interface{}, error) {
	token, err := jwt.ParseWithClaims(bearToken, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if claims, ok := token.Claims.(*Claim); ok && token.Valid {
		db := bd.ConnectDB()
		sqldb, _ := db.DB()
		defer sqldb.Close()

		result := db.Table("usuarios_jwt").Select("count(id)").Where("id_user = ? and token = ? and validation = 'true'", claims.User.IdUser, bearToken).Scan(&dato)

		if dato.Count == 0 {
			return "", fmt.Errorf("%s", "Invalid Token")
		}

		return claims.Usuario.Id, nil
	}
	return "", err
}

// funcion para extraer el id del cliente del token
// se retorna el id del usuario, el token, el codigo del error y el error
func TokenVerification(token string) (string, string, int, error) {
	if token == "" {
		return "", "", 106, fmt.Errorf("token is empty")
	}
	bearToken := strings.ReplaceAll(token, " ", "")
	//se valida el token y se extrae el id del usuario del mismo
	idUserJwt, err := validateJWT(bearToken)
	if err != nil {
		return "", "", 106, err

	}

	//se parsea el id del usuario
	idUser := fmt.Sprintf("%v", idUserJwt)
	return idUser, bearToken, 0, nil
}

// Se retorna si el token es válido, la fecha y hora de expiración y si hay error
func CheckTokenValidityDate(bearToken string) (bool, time.Time, error) {
	if bearToken == "" {
		return false, time.Time{}, fmt.Errorf("token is empty")
	}
	bearToken = strings.ReplaceAll(bearToken, " ", "")
	token, err := jwt.ParseWithClaims(bearToken, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return false, time.Time{}, err
	}

	if claims, ok := token.Claims.(*Claim); ok && token.Valid {
		return token.Valid, claims.ExpiresAt.Time, nil
	} else {
		return false, time.Time{}, fmt.Errorf("token expired or not valid")
	}
}
