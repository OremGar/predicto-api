package controladores

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/OremGar/predicto-api/bd"
	"github.com/OremGar/predicto-api/respuestas"
	"gorm.io/gorm"
)

func CambiarNombreAnalicto(w http.ResponseWriter, r *http.Request) {
	var id int
	var nuevoNombre string
	var err error
	var existe bool

	var db *gorm.DB = bd.ConnectDB()
	sqldb, _ := db.DB()
	defer sqldb.Close()

	id, err = strconv.Atoi(r.FormValue("id"))
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("id de analicto en formato incorrecto"))
		return
	}

	nuevoNombre = r.FormValue("nombre")
	if nuevoNombre == "" {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("el nombre del analicto no debe de estar vacío"))
		return
	}

	fmt.Println(nuevoNombre)

	result := db.Raw("SELECT COUNT(*)>0 FROM analictos WHERE id = ?", id).Scan(&existe)
	if result.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("internal error: %v", result.Error))
		return
	}

	fmt.Println(existe)

	if !existe {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("el analicto no existe"))
		return
	}

	result = db.Raw("UPDATE analictos SET nombre = ? WHERE id = ?", nuevoNombre, id)
	if result.Error != nil {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("el analicto no existe"))
		return
	}

	fmt.Println(result.RowsAffected)

	respuestas.JsonResponse(w, http.StatusOK, nil, 0, nil)
}
