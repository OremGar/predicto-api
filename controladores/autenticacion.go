package controladores

import (
	"net/http"

	"github.com/OremGar/predicto-api/respuestas"
)

func Prueba(w http.ResponseWriter, r *http.Request) {
	type respuesta struct {
		Saludo string
	}
	respuestas.JsonResponse(w, http.StatusOK, respuesta{Saludo: "Todo ok"}, 0, nil)
}
