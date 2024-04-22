package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/OremGar/predicto-api/bd"
	"github.com/OremGar/predicto-api/configuraciones"
	"github.com/OremGar/predicto-api/modelos"
	"github.com/OremGar/predicto-api/respuestas"
	"gorm.io/gorm"
)

// Middleware para válidar usuarios y tokens
func ValidarToken(peticion http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string = r.Header.Get("Authorization") //Se obtiene el token del encabezado
		var usuario modelos.Usuarios                     //Objeto usuario
		var existe bool                                  //Booleano para válidar existencia
		var err error

		var db *gorm.DB
		db, err = bd.ConnectDB()
		if err != nil {
			respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error en la bd: %v", err))
			return
		}
		sqldb, _ := db.DB()
		defer sqldb.Close()

		if strings.HasPrefix(r.URL.Path, "/api/v1/cuenta") || strings.HasPrefix(r.URL.Path, "/api/v1/general") { //Si el endpoint a consultar tiene la ruta /api/v1/cuenta, no es necesario hacer el resto de validaciones
			peticion.ServeHTTP(w, r)
			return
		}

		_, claims, err := configuraciones.ValidarJWT(token) //Se válida el JWT
		if err != nil {
			respuestas.SetError(w, http.StatusUnauthorized, 100, err)
			return
		}

		usuario.Id = claims.IdUsuario //Se obtiene el id del usuario de los claims del JWT

		result := db.First(&usuario) //Se busca el usuario por su ID
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				respuestas.SetError(w, http.StatusUnauthorized, 100, fmt.Errorf("el usuario no existe")) //Si no existe usuario, se retorna un error
				return
			}
			respuestas.SetError(w, http.StatusUnauthorized, 100, result.Error)
			return
		}

		result = db.Model(&modelos.UsuariosJwt{}).Select("count(*) > 0").Where("token = ?", token).Find(&existe) //Se busca la existencia del JWT
		if result.Error != nil {
			respuestas.SetError(w, http.StatusInternalServerError, 100, result.Error)
			return
		}
		if !existe {
			respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("el token no existe en la bd")) //Si no existe significa que nunca se generó aquí, por lo tanto no es un token válido
			return
		}

		//Si hay tokens más nuevos que el de la consulta, entonces no se permite el acceso porque significa que ya se inició sesión en otros dispositivos
		result = db.Raw("SELECT count(*) > 0 FROM usuarios_jwt WHERE id_usuario = ? AND fecha_inicio > (SELECT fecha_inicio FROM usuarios_jwt WHERE token = ? LIMIT 1)", usuario.Id, token).Find(&existe)
		if result.Error != nil {
			respuestas.SetError(w, http.StatusInternalServerError, 100, result.Error)
			return
		}
		if existe {
			respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("token no válido, existen tokens más nuevos"))
			return
		}

		peticion.ServeHTTP(w, r)
	})
}
