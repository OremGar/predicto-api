package modelos

import (
	"fmt"
	"time"

	"github.com/OremGar/predicto-api/bd"
	"gorm.io/gorm"
)

const (
	ROTOR_INTERNO string = "interno"
	ROTOR_EXTERNO string = "externo"
)

type Motores struct {
	Id       int    `json:"id"`
	Potencia int    `json:"potencia"`
	Rotor    string `json:"rotor"`
}

type MotoresVibraciones struct {
	Id      int       `json:"id"`
	IdMotor int       `json:"id_motor"`
	Hora    time.Time `json:"hora"`
	EjeX    float64   `json:"eje_x"`
	EjeY    float64   `json:"eje_y"`
	EjeZ    float64   `json:"eje_z"`
}

type MotoresEstados struct {
	Id       int       `json:"id"`
	IdMotor  int       `json:"id_motor"`
	Fecha    time.Time `json:"fecha"`
	Estado   string    `json:"estado"`
	Contador string    `json:"-"`
}

func MotorExiste(id int) (Motores, error) {
	var motor Motores = Motores{}
	var err error

	var db *gorm.DB
	db, err = bd.ConnectDB()
	if err != nil {
		return Motores{}, fmt.Errorf("error en la bd: %v", err)
	}
	sqldb, _ := db.DB()
	defer sqldb.Close()

	resultado := db.Model(&Motores{}).Where("id = ?", id).First(&motor)
	if resultado.Error != nil {
		if resultado.Error == gorm.ErrRecordNotFound {
			return Motores{}, fmt.Errorf("el motor no existe")
		}
		return Motores{}, fmt.Errorf("error interno: %v", resultado.Error)
	}

	return motor, nil
}
