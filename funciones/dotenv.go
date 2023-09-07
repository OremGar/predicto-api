package funciones

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func GetDotEnvVar(key string) string {
	//Se obtiene el directorio actual
	ruta, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("Error al obtener la ruta actual: %v\n", err)
	}

	godotenv.Load()

	// load .env file
	err = godotenv.Load(ruta + "/env")

	if err != nil {
		return ""
	}

	return os.Getenv(key)
}
