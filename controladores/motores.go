package controladores

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/OremGar/predicto-api/bd"
	"github.com/OremGar/predicto-api/modelos"
	"github.com/OremGar/predicto-api/respuestas"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func ObtieneMotores(w http.ResponseWriter, r *http.Request) {
	type respuestaStruct struct {
		Motor       modelos.Motores
		FechaInicio time.Time
		FechaFin    time.Time
	}

	var motores []modelos.Motores = []modelos.Motores{}
	var err error
	var respuesta []respuestaStruct

	var db *gorm.DB
	db, err = bd.ConnectDB()
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error en la bd: %v", err))
		return
	}
	sqldb, _ := db.DB()
	defer sqldb.Close()

	resultado := db.Model(&modelos.Motores{}).Find(&motores)
	if resultado.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error al obtener los motores: %v", resultado.Error))
		return
	}

	for _, motor := range motores {
		var elemento respuestaStruct = respuestaStruct{}
		elemento.Motor = motor

		result := db.Raw("SELECT hora FROM motores_vibraciones WHERE id_motor = ? ORDER BY hora ASC LIMIT 1", motor.Id).Scan(&elemento.FechaInicio)
		if result.Error != nil {
			respuestas.SetError(w, http.StatusInternalServerError, 102, fmt.Errorf("error buscando la fecha del primer paquete: %v", result.Error))
			return
		}

		result = db.Raw("SELECT hora FROM motores_vibraciones WHERE id_motor = ? ORDER BY hora DESC LIMIT 1", motor.Id).Scan(&elemento.FechaFin)
		if result.Error != nil {
			respuestas.SetError(w, http.StatusInternalServerError, 103, fmt.Errorf("error buscando la fecha del segundo paquete: %v", result.Error))
			return
		}

		respuesta = append(respuesta, elemento)
	}

	respuestas.JsonResponse(w, http.StatusOK, respuesta, 0, nil)
}

func ObtieneVibracionPeriodo(w http.ResponseWriter, r *http.Request) {
	type Respuesta struct {
		FechaInicio time.Time
		FechaFin    time.Time
	}

	var vars map[string]string = mux.Vars(r)

	var motor modelos.Motores = modelos.Motores{}
	var err error

	var respuesta Respuesta = Respuesta{}

	var db *gorm.DB
	db, err = bd.ConnectDB()
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error en la bd: %v", err))
		return
	}
	sqldb, _ := db.DB()
	defer sqldb.Close()

	motor.Id, err = strconv.Atoi(vars["id"])
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("el id del motor no está en el formato correcto"))
		return
	}

	_, err = modelos.MotorExiste(motor.Id)
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 101, err)
		return
	}

	result := db.Raw("SELECT hora FROM motores_vibraciones WHERE id_motor = ? ORDER BY hora ASC LIMIT 1", motor.Id).Scan(&respuesta.FechaInicio)
	if result.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 102, fmt.Errorf("error buscando la fecha del primer paquete: %v", result.Error))
		return
	}

	result = db.Raw("SELECT hora FROM motores_vibraciones WHERE id_motor = ? ORDER BY hora DESC LIMIT 1", motor.Id).Scan(&respuesta.FechaFin)
	if result.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 103, fmt.Errorf("error buscando la fecha del segundo paquete: %v", result.Error))
		return
	}

	respuestas.JsonResponse(w, http.StatusOK, respuesta, 0, nil)
}

func ObtieneVibracionesMotores(w http.ResponseWriter, r *http.Request) {
	var motor modelos.Motores = modelos.Motores{}
	var vibraciones []modelos.MotoresVibraciones = []modelos.MotoresVibraciones{}

	var fecInicio time.Time
	var fecFinal time.Time

	var err error

	var db *gorm.DB
	db, err = bd.ConnectDB()
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error en la bd: %v", err))
		return
	}
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

	fecInicio, err = time.Parse("2006-01-02", r.FormValue("fecInicio"))
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("la fecha de inicio no está en el formato correcto YYYY-MM-DD"))
		return
	}

	fecFinal, err = time.Parse("2006-01-02", r.FormValue("fecFin"))
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
