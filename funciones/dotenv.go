package funciones

import (
	"os"

	"github.com/joho/godotenv"
)

func GetDotEnvVar(key string) string {
	godotenv.Load()

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		return ""
	}

	return os.Getenv(key)
}
