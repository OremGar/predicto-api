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
	r.HandleFunc("/api/v1/cuenta/RecuperaContrasena", controladores.RecuperaContrasena).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/cuenta/NuevaContrasena", controladores.NuevaContrasena).Methods(http.MethodPut)
	r.HandleFunc("/api/v1/cuenta/ValidaOTPContrasena/{codigo:[0-9]+}", controladores.ValidaOTPNvaContrasena).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/cuenta/ValidaOTPLogin", controladores.ValidaOtpLogin).Methods(http.MethodPost)

	//Administrador

	r.HandleFunc("/api/v1/prueba/saludo", controladores.Prueba).Methods(http.MethodGet)

	//Motores
	r.HandleFunc("/api/v1/motores", controladores.ObtieneMotores).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/motores/vibraciones", controladores.ObtieneVibracionesMotores).Methods(http.MethodPatch)

	wrappedMux = middlewares.ValidarToken(r.ServeHTTP)
	wrappedMux = middlewares.ValidarAdmin(wrappedMux)

	corsOpc = cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://predicto.ddns.net:8080",
			"http://predicto.ddns.net:8081",
			"https://predicto.ddns.net:8080",
			"https://predicto.ddns.net:8081",
			"https://predicto.ddns.net:8083",
			"https://predicto.ddns.net:8083",
			"http://192.168.1.25:8080",
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

	if PUERTO == "" {
		PUERTO = "8081"
	}

	log.Println("Iniciando servidor en el puerto:", PUERTO)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PUERTO), corsOpt.Handler(r)))
}
