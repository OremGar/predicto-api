package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/OremGar/predicto-api/controladores"
	"github.com/OremGar/predicto-api/funciones"
	"github.com/OremGar/predicto-api/middlewares"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	PUERTO = funciones.GetDotEnvVar("PUERTO")
)

func Router() (http.Handler, *cors.Cors) {
	var r *mux.Router = mux.NewRouter() //Creación del router
	var corsOpc *cors.Cors
	var wrappedMux http.Handler

	r.HandleFunc("/api/v1/cuenta/SignUp", controladores.SignUp).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/cuenta/SignIn", controladores.SignIn).Methods(http.MethodPost)

	wrappedMux = middlewares.ValidarToken(r.ServeHTTP)

	corsOpc = cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://oremserver.duckdns.org:8081",
			"http://localhost:3000",
			"*",
		},

		AllowedMethods: []string{
			"*",
			http.MethodGet, //http methods for your app
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},

		AllowedHeaders: []string{
			"*",
			"Content-Type",
			"Authorization",
		},

		AllowCredentials: true,
	})

	return wrappedMux, corsOpc
}

func main() {
	r, corsOpt := Router()

	log.Println("Iniciando servidor en el puerto:", PUERTO)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PUERTO), corsOpt.Handler(r)))
}
