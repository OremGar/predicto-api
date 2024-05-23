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

func ObtenerTolerancias(w http.ResponseWriter, r *http.Request) {
	var err error

	var vars map[string]string

	var tolerancia modelos.Tolerancia
	var motor modelos.Motores

	var db *gorm.DB
	var result *gorm.DB

	db, err = bd.ConnectDB()
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error interno"))
		return
	}

	vars = mux.Vars(r)

	motor.Id, err = strconv.Atoi(vars["id"])
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("id de motor no esta en el formato correcto"))
		return
	}

	motor, err = modelos.MotorExiste(motor.Id)
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, err)
		return
	}

	result = db.Model(&modelos.Tolerancia{}).Where("id_motor = ?", motor.Id).First(&tolerancia)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("no existe tolerancia para este motor"))
			return
		}
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("internal error"))
		return
	}

	respuestas.JsonResponse(w, http.StatusOK, tolerancia, 0, nil)
}

func ActualizaTolerancia(w http.ResponseWriter, r *http.Request) {
	var err error

	var vars map[string]string

	var tolerancia modelos.Tolerancia
	var motor modelos.Motores

	var db *gorm.DB
	var result *gorm.DB

	db, err = bd.ConnectDB()
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error interno"))
		return
	}

	vars = mux.Vars(r)

	motor.Id, err = strconv.Atoi(vars["id"])
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("id de motor no esta en el formato correcto"))
		return
	}

	motor, err = modelos.MotorExiste(motor.Id)
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, err)
		return
	}

	result = db.Model(&modelos.Tolerancia{}).Where("id_motor = ?", motor.Id).First(&tolerancia)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("no existe tolerancia para este motor"))
			return
		}
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("internal error"))
		return
	}

	switch r.FormValue("roll") {
	case "1":
		tolerancia.RollMax = modelos.GRAVITACION_SENSIBILIDAD_BAJA
		tolerancia.RollMin = -modelos.GRAVITACION_SENSIBILIDAD_BAJA

	case "2":
		tolerancia.RollMax = modelos.GRAVITACION_SENSIBILIDAD_MEDIA
		tolerancia.RollMin = -modelos.GRAVITACION_SENSIBILIDAD_MEDIA

	case "3":
		tolerancia.RollMax = modelos.GRAVITACION_SENSIBILIDAD_ALTA
		tolerancia.RollMin = -modelos.GRAVITACION_SENSIBILIDAD_ALTA

	default:
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("roll: valor incorrecto"))
	}

	switch r.FormValue("pitch") {
	case "1":
		tolerancia.PitchMax = modelos.GRAVITACION_SENSIBILIDAD_BAJA
		tolerancia.PitchMin = -modelos.GRAVITACION_SENSIBILIDAD_BAJA

	case "2":
		tolerancia.PitchMax = modelos.GRAVITACION_SENSIBILIDAD_MEDIA
		tolerancia.PitchMin = -modelos.GRAVITACION_SENSIBILIDAD_MEDIA

	case "3":
		tolerancia.PitchMax = modelos.GRAVITACION_SENSIBILIDAD_ALTA
		tolerancia.PitchMin = -modelos.GRAVITACION_SENSIBILIDAD_ALTA

	default:
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("pitch: valor incorrecto"))
	}

	if r.FormValue("temp") != "" {
		var tmp float64
		tmp, err = strconv.ParseFloat(r.FormValue("temp"), 32)
		if err != nil {
			respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("temperatura no está en el formato correcto"))
			return
		}

		tolerancia.TempMax = float32(tmp)
	} else {
		tolerancia.TempMax = 90
	}

	result = db.Save(&tolerancia)
	if result.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("internal error"))
		return
	}

	respuestas.JsonResponse(w, http.StatusOK, tolerancia, 0, nil)
}
