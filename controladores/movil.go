package controladores

import (
	"fmt"
	"net/http"

	"github.com/OremGar/predicto-api/bd"
	"github.com/OremGar/predicto-api/configuraciones"
	"github.com/OremGar/predicto-api/modelos"
	"github.com/OremGar/predicto-api/respuestas"
	"gorm.io/gorm"
)

func GuardaTokenFirebase(w http.ResponseWriter, r *http.Request) {
	var err error
	var db *gorm.DB

	var token string
	var jwt string
	var claims modelos.Claims
	var tokenFirebase modelos.TokenFirebase

	db, err = bd.ConnectDB()
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error interno"))
		return
	}

	sqldb, _ := db.DB()
	defer sqldb.Close()

	token = r.FormValue("token")
	if token == "" {
		respuestas.SetError(w, http.StatusBadRequest, 101, fmt.Errorf("el token esta vacio"))
	}

	jwt = r.Header.Get("Authorization")

	_, claims, err = configuraciones.ValidarJWT(jwt)
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 102, err)
		return
	}

	_, err = modelos.ChecarSiUsuarioExiste(claims.IdUsuario)
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 103, err)
		return
	}

	tokenFirebase.IdUsuario = claims.IdUsuario
	tokenFirebase.Token = token

	result := db.Create(&tokenFirebase)
	if result.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 104, fmt.Errorf("error interno"))
		return
	}

	respuestas.JsonResponse(w, http.StatusOK, nil, 0, nil)
}
