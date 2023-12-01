package controladores

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/OremGar/predicto-api/respuestas"
)

func CambiarNombreAnalicto(w http.ResponseWriter, r *http.Request) {
	var id int
	var err error

	id, err = strconv.Atoi(r.FormValue("id"))
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("id de analicto en formato incorrecto"))
		return
	}

	fmt.Println(id)

	respuestas.JsonResponse(w, http.StatusOK, nil, 0, nil)
}
