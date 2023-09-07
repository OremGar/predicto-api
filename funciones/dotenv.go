package funciones

import (
	"os"

	"github.com/joho/godotenv"
)

func GetDotEnvVar(key string) string {
	//Se obtiene el directorio actual
	var ruta = "~/Predicto/predicto-api"

	godotenv.Load()

	// load .env file
	err := godotenv.Load(ruta + "/env")

	if err != nil {
		return ""
	}

	return os.Getenv(key)
}
