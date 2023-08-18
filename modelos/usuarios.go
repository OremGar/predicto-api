package modelos

import (
	"fmt"
	"time"

	"github.com/OremGar/predicto-api/bd"
	"gorm.io/gorm"
)

// Modelos
type Usuarios struct {
	Id         int
	Nombre     string
	Apellidos  string
	Usuario    string
	Correo     string
	Contrasena string
	Telefono   string
}

// Métodos
type UsuariosJwt struct {
	IdUsuario   int
	Token       string
	FechaInicio time.Time
}

func ChecarSiUsuarioExiste(id int) (Usuarios, error) {
	var usuario Usuarios = Usuarios{}

	var db *gorm.DB = bd.ConnectDB()
	sqldb, _ := db.DB()
	defer sqldb.Close()

	result := db.Model(&usuario).Where("id = ?", id).First(&usuario)
	if result.Error != nil {
		return Usuarios{}, result.Error
	}
	if result.RowsAffected <= 0 {
		return Usuarios{}, fmt.Errorf("el usuario no existe")
	}

	return usuario, nil
}

// Función para validar que el objeto usuario no le falte nada
func ValidarInfoUsuarios(usuario *Usuarios) error {
	if usuario.Nombre == "" {
		return fmt.Errorf("falta el nombre")
	}
	if usuario.Apellidos == "" {
		return fmt.Errorf("faltan los apellidos")
	}
	if usuario.Usuario == "" {
		return fmt.Errorf("falta el usuario")
	}
	if usuario.Correo == "" {
		return fmt.Errorf("falta el correo")
	}
	if usuario.Contrasena == "" {
		return fmt.Errorf("falta la contrasena")
	}
	if usuario.Telefono == "" {
		return fmt.Errorf("falta el telefono")
	}

	return nil
}
