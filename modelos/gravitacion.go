package modelos

import "time"

type Gravitacion struct {
	Id          int       `json:"id"`
	IdMotor     int       `json:"id_motor"`
	Roll        float64   `json:"roll"`
	Pitch       float64   `json:"pitch"`
	Fecha       time.Time `json:"fecha"`
	Temperatura float64   `json:"temperatura"`
}

type Anomalias struct {
	Id            int    `json:"int"`
	IdGravitacion int    `json:"id_gravitacion"`
	Anomalia      string `json:"anomalia"`
}

func (Gravitacion) TableName() string {
	return "gravitacion"
}

func (Anomalias) TableName() string {
	return "anomalias"
}
