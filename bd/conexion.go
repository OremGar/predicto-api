package bd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/OremGar/predicto-api/funciones"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Se obtienen las variables del archivo .env
var (
	host       = funciones.GetDotEnvVar("HOST")
	usuario    = funciones.GetDotEnvVar("USUARIO")
	contrasena = funciones.GetDotEnvVar("CONTRASENA")
	base_datos = funciones.GetDotEnvVar("BASE_DATOS")
	puerto, _  = strconv.Atoi(funciones.GetDotEnvVar("PUERTO_BD"))
)

// Función que realiza una conexión a la BD y retorna un objeto para realizar las operaciones con ella
func ConnectDB() *gorm.DB {
	//Connect to DB
	var DB *gorm.DB
	var dsn string = fmt.Sprintf("host=%s user=%s password=%s  dbname=%s port=%d  sslmode=disable", host, usuario, contrasena, base_datos, puerto)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("Error en la conexión a la BD %v", err)
		return nil
	}

	db, err := DB.DB()

	if err := db.Ping(); err != nil {
		log.Fatalln("Error haciendo ping en la BD  " + err.Error())
		return nil
	}

	db.SetConnMaxIdleTime(time.Minute * 5)
	//Se validan las conexiones a la BD
	if err != nil {
		fmt.Printf("Error en la conexión a la BD %v", err)
		return nil
	}
	if DB.Error != nil {
		fmt.Printf("Cualquier error en la conexión a la BD %v" + err.Error())
		return nil
	}
	log.Println("Conexión a BD exitosa")
	return DB
}
