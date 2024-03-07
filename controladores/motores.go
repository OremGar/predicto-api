package controladores

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/OremGar/predicto-api/bd"
	"github.com/OremGar/predicto-api/modelos"
	"github.com/OremGar/predicto-api/respuestas"
	"gorm.io/gorm"
)

func ObtieneMotores(w http.ResponseWriter, r *http.Request) {
	var motores []modelos.Motores = []modelos.Motores{}

	var db *gorm.DB = bd.ConnectDB()
	sqldb, _ := db.DB()
	defer sqldb.Close()

	resultado := db.Model(&modelos.Motores{}).Select(&motores)
	if resultado.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error al obtener los motores: %v", resultado.Error))
		return
	}

	respuestas.JsonResponse(w, http.StatusOK, motores, 0, nil)
}

func ObtieneVibracionesMotores(w http.ResponseWriter, r *http.Request) {
	var motor modelos.Motores = modelos.Motores{}
	var vibraciones modelos.MotoresVibraciones = modelos.MotoresVibraciones{}

	var fecInicio time.Time
	var fecFinal time.Time

	var err error

	var db *gorm.DB = bd.ConnectDB()
	sqldb, _ := db.DB()
	defer sqldb.Close()

	motor.Id, err = strconv.Atoi(r.FormValue("id"))
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("el id del motor no está en el formato correcto"))
		return
	}

	motor, err = modelos.MotorExiste(motor.Id)
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 100, err)
		return
	}

	fecInicio, err = time.Parse("2006-01-02", r.FormValue("fec_inicio"))
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("la fecha de inicio no está en el formato correcto YYYY-MM-DD"))
		return
	}

	fecFinal, err = time.Parse("2006-01-02", r.FormValue("fec_fin"))
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("la fecha de fin no está en el formato correcto YYYY-MM-DD"))
		return
	}

	resultado := db.Model(&modelos.MotoresVibraciones{}).Order("hora ASC").Where("id_motor = ? AND hora >= ? AND hora <= ?", motor.Id, fecInicio, fecFinal).Find(&vibraciones)
	if resultado.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error buscando las vibraciones: %v", resultado.Error))
		return
	}

	respuestas.JsonResponse(w, http.StatusOK, vibraciones, 0, nil)
}
