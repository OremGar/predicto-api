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

func ObtieneAnomalias(w http.ResponseWriter, r *http.Request) {
	var err error

	var vars map[string]string

	var db *gorm.DB
	var result *gorm.DB

	var motor modelos.Motores
	var anomalias []modelos.Anomalias

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

	subQuery1 := db.Model(&modelos.Gravitacion{}).Order("fecha DESC").Select("id").Where("id_motor = ?", motor.Id)
	result = db.Model(&modelos.Anomalias{}).Order("id DESC").Where("id_gravitacion IN (?)", subQuery1).Find(&anomalias)
	if result.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error interno"))
		return
	}

	respuestas.JsonResponse(w, http.StatusOK, anomalias, 0, nil)
}
