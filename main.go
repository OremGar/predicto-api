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
	r.HandleFunc("/api/v1/motores/vibraciones/periodo/{id:[0-9]+}", controladores.ObtieneVibracionPeriodo).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/motores/vibraciones", controladores.ObtieneVibracionesMotores).Methods(http.MethodPatch)
	r.HandleFunc("/api/v1/motores/estados/{id:[0-9]+}", controladores.ObtieneEstados).Methods(http.MethodGet)

	//Tolerancia
	r.HandleFunc("/api/v1/motores/tolerancia/{id:[0-9]+}", controladores.ObtenerTolerancias).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/motores/tolerancia/{id:[0-9]+}", controladores.ActualizaTolerancia).Methods(http.MethodPut)

	//Anomalías
	r.HandleFunc("/api/v1/motores/anomalias/{id:[0-9]+}", controladores.ObtieneAnomalias).Methods(http.MethodGet)

	//Gravitaciones
	r.HandleFunc("/api/v1/motores/gravitaciones/{id:[0-9]+}", controladores.ObtieneGravitaciones).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/motores/gravitacion/{id:[0-9]+}", controladores.ObtieneGravitacion).Methods(http.MethodGet)

	//Firebase
	r.HandleFunc("/api/v1/GuardaTokenFB", controladores.GuardaTokenFirebase).Methods(http.MethodPost)

	//Temporal
	r.HandleFunc("/api/v1/general", controladores.ResetearBD).Methods(http.MethodGet)

	wrappedMux = middlewares.ValidarToken(r.ServeHTTP)
	wrappedMux = middlewares.ValidarAdmin(wrappedMux)

	//CORS
	corsOpc = cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://predicto.ddns.net",
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
			http.MethodPatch,
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
