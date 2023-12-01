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

func ValidarAdmin(peticion http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string = r.Header.Get("Authorization")
		var usuario modelos.Usuarios
		var existe bool

		var db *gorm.DB = bd.ConnectDB()
		sqldb, _ := db.DB()
		defer sqldb.Close()

		fmt.Println(r.URL.Path)
		if !strings.HasPrefix(r.URL.Path, "/api/v1/cuenta/admin") {
			peticion.ServeHTTP(w, r)
			return
		}

		_, claims, err := configuraciones.ValidarJWT(token)
		if err != nil {
			respuestas.SetError(w, http.StatusUnauthorized, 100, fmt.Errorf("token is not valid"))
			return
		}

		usuario.Id = claims.IdUsuario

		result := db.First(&usuario)
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
		result = db.Raw("SELECT count(*) > 0 FROM usuarios_jwt WHERE id_usuario = ? AND fecha_inicio > (SELECT fecha_inicio FROM usuarios_jwt WHERE token = ?)", usuario.Id, token).Find(&existe)
		if result.Error != nil {
			respuestas.SetError(w, http.StatusInternalServerError, 100, result.Error)
			return
		}
		if existe {
			respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("token no válido, existen tokens más nuevos"))
			return
		}

		if usuario.TipoUsuario != modelos.TIPO_USUARIO_ADMIN && usuario.TipoUsuario != modelos.TIPO_USUARIO_USUARIO_ADMINISTRADOR {
			respuestas.SetError(w, http.StatusUnauthorized, 100, fmt.Errorf("el usuario no está autorizado para acceder a este recurso"))
			return
		}

		peticion.ServeHTTP(w, r)
	})
}
