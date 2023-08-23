package controladores

import (
	"fmt"
	"net/http"

	"github.com/OremGar/predicto-api/bd"
	"github.com/OremGar/predicto-api/funciones"
	"github.com/OremGar/predicto-api/modelos"
	"github.com/OremGar/predicto-api/respuestas"
	"gorm.io/gorm"
)

func SignIn(w http.ResponseWriter, r *http.Request) {

}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var usuario modelos.Usuarios = modelos.Usuarios{}

	var db *gorm.DB = bd.ConnectDB()
	sqldb, _ := db.DB()
	defer sqldb.Close()

	usuario.Nombre = r.FormValue("nombre")
	usuario.Apellidos = r.FormValue("apellidos")
	usuario.Correo = r.FormValue("correo")
	if !funciones.ValidaCorreo(usuario.Correo) {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("el correo no está en el formato correcto"))
		return
	}
	usuario.Usuario = r.FormValue("usuario")
	contrasena, err := funciones.HashContrasena(r.FormValue("contrasena"))
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, err)
		return
	}
	usuario.Contrasena = contrasena
	usuario.Telefono = r.FormValue("telefono")

	err = modelos.ValidarInfoUsuarios(&usuario)
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 100, err)
		return
	}

	result := db.Save(&usuario).First(&usuario)
	if result.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, result.Error)
		return
	}

	respuestas.JsonResponse(w, http.StatusCreated, usuario.Id, 0, nil)
}

func Saludo(w http.ResponseWriter, r *http.Request) {
	respuestas.JsonResponse(w, http.StatusOK, "Saludos", 0, nil)
}
