package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/OremGar/predicto-api/controladores"
	"github.com/OremGar/predicto-api/funciones"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	PUERTO = funciones.GetDotEnvVar("PUERTO")
)

func Router() (http.Handler, *cors.Cors) {
	var r *mux.Router = mux.NewRouter() //Creación del router
	var corsOpc *cors.Cors

	r.HandleFunc("/api/v1/signup", controladores.SignUp).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/signin", controladores.SignIn).Methods(http.MethodPost)

	corsOpc = cors.New(cors.Options{
		AllowedOrigins: []string{
			"*",
		},

		AllowedMethods: []string{
			"*",
		},

		AllowedHeaders: []string{
			"*",
			"Content-Type",
			"Authorization",
		},

		AllowCredentials: true,
	})

	return r, corsOpc
}

func main() {
	r, corsOpt := Router()

	log.Println("Iniciando servidor en el puerto:", PUERTO)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PUERTO), corsOpt.Handler(r)))
}
