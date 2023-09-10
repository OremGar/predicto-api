package modelos

import "bytes"

type Correo struct {
	Origen  string
	Destino string
	Asunto  string
	Cuerpo  bytes.Buffer
}
