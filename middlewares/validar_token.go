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

func ValidarToken(peticion http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string = r.Header.Get("Authorization")
		var usuario modelos.Usuarios

		var db *gorm.DB = bd.ConnectDB()
		sqldb, _ := db.DB()
		defer sqldb.Close()

		if strings.HasPrefix(r.URL.Path, "/api/v1/cuenta") {
			peticion.ServeHTTP(w, r)
			return
		}

		_, claims, err := configuraciones.ValidarJWT(token)
		if err != nil {
			respuestas.SetError(w, http.StatusUnauthorized, 100, fmt.Errorf("token no válido"))
			return
		}

		usuario.Id = claims.IdUsuario

		result := db.First(&usuario)
		if result.Error != nil {
			respuestas.SetError(w, http.StatusUnauthorized, 100, result.Error)
			return
		}
		if result.RowsAffected == 0 {
			respuestas.SetError(w, http.StatusUnauthorized, 100, fmt.Errorf("el usuario no existe"))
			return
		}

		result = db.Model(&modelos.UsuariosJwt{}).Where("token = ?", token)
		if result.Error != nil {
			respuestas.SetError(w, http.StatusInternalServerError, 100, result.Error)
			return
		}
		if result.RowsAffected == 0 {
			respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("el token no existe en la bd"))
			return
		}

		result = db.Raw("SELECT * FROM usuarios_jwt WHERE id_usuario = ? AND fecha_inicio > (SELECT fecha_inicio FROM usuarios_jwt WHERE token = ?)", usuario.Id, token)
		if result.Error != nil {
			respuestas.SetError(w, http.StatusInternalServerError, 100, result.Error)
			return
		}
		if result.RowsAffected > 0 {
			respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("token no válido, existen tokens más nuevos"))
			return
		}

		peticion.ServeHTTP(w, r)
	})
}
