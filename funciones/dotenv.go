package funciones

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetDotEnvVar(key string) string {
	//Se obtiene el directorio actual
	ruta, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error al obtener la ruta actual: %v\n", err)
	}

	godotenv.Load()

	// load .env file
	err = godotenv.Load(ruta + "/.env")

	if err != nil {
		return ""
	}

	return os.Getenv(key)
}
