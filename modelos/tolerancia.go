package modelos

const (
	GRAVITACION_SENSIBILIDAD_BAJA  float32 = 2.5
	GRAVITACION_SENSIBILIDAD_MEDIA float32 = 2
	GRAVITACION_SENSIBILIDAD_ALTA  float32 = 1
)

type Tolerancia struct {
	Id       int     `json:"int"`
	IdMotor  int     `json:"id_motor"`
	RollMax  float32 `json:"roll_max"`
	RollMin  float32 `json:"roll_min"`
	TempMax  float32 `json:"temp_max"`
	PitchMax float32 `json:"pitch_max"`
	PitchMin float32 `json:"pitch_min"`
}

func (Tolerancia) TableName() string {
	return "tolerancia"
}
