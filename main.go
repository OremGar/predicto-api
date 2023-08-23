package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/OremGar/predicto-api/bd"
	"github.com/OremGar/predicto-api/funciones"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

var (
	PUERTO = funciones.GetDotEnvVar("PUERTO")
)

func Router() (http.Handler, *cors.Cors) {
	var r *mux.Router = mux.NewRouter() //Creación del router
	var corsOpc *cors.Cors

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

	var db *gorm.DB = bd.ConnectDB()
	sqldb, _ := db.DB()
	defer sqldb.Close()

	log.Println("Iniciando servidor en el puerto:", PUERTO)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PUERTO), corsOpt.Handler(r)))
}
