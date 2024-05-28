package controladores

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/OremGar/predicto-api/bd"
	"github.com/OremGar/predicto-api/modelos"
	"github.com/OremGar/predicto-api/respuestas"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func ObtieneGravitaciones(w http.ResponseWriter, r *http.Request) {
	var err error

	var vars map[string]string

	var db *gorm.DB
	var result *gorm.DB

	var motor modelos.Motores
	var gravitaciones []modelos.Gravitacion

	vars = mux.Vars(r)

	db, err = bd.ConnectDB()
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error interno"))
		return
	}

	motor.Id, err = strconv.Atoi(vars["id"])
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 101, fmt.Errorf("el id del motor no está en el formato correcto"))
		return
	}

	motor, err = modelos.MotorExiste(motor.Id)
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 102, err)
		return
	}

	result = db.Model(&modelos.Gravitacion{}).Order("fecha DESC").Where("id_motor = ?", motor.Id).Find(&gravitaciones)
	if result.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error interno"))
		return
	}

	respuestas.JsonResponse(w, http.StatusOK, gravitaciones, 0, nil)
}

func ObtieneGravitacion(w http.ResponseWriter, r *http.Request) {
	var err error

	var vars map[string]string

	var db *gorm.DB
	var result *gorm.DB

	var gravitacion modelos.Gravitacion

	vars = mux.Vars(r)

	db, err = bd.ConnectDB()
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error interno"))
		return
	}

	gravitacion.Id, err = strconv.Atoi(vars["id"])
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 101, fmt.Errorf("el id de la gravitación no está en el formato correcto"))
		return
	}

	result = db.Model(&modelos.Gravitacion{}).Where("id = ?", gravitacion.Id).Limit(1).Find(&gravitacion)
	if result.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error interno"))
		return
	}

	if result.RowsAffected == 0 {
		respuestas.JsonResponse(w, http.StatusOK, nil, 0, nil)
		return
	}
	respuestas.JsonResponse(w, http.StatusOK, gravitacion, 0, nil)
}
