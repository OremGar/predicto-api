package respuestas

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Errors Error       `json:"errors"`
}

type Error struct {
	Code       int      `json:"code"`
	Decription []string `json:"error"`
}

var MetaErrors []string

func JsonResponse(w http.ResponseWriter, status int, data interface{}, code int, errors []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Status: status,
		Data:   data,
		Errors: Error{
			Code:       code,
			Decription: errors,
		},
	})
}

func SetError(w http.ResponseWriter, status int, code int, err error) {
	MetaErrors = nil
	MetaErrors = append(MetaErrors, err.Error())
	JsonResponse(w, status, nil, code, MetaErrors)
}
