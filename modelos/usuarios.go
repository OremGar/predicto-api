package modelos

import (
	"fmt"
	"time"

	"github.com/OremGar/predicto-api/bd"
	"gorm.io/gorm"
)

// Modelos
type Usuarios struct {
	Id         int    `json:"id"`
	Nombre     string `json:"nombre"`
	Apellidos  string `json:"apellidos"`
	Usuario    string `json:"usuario"`
	Correo     string `json:"correo"`
	Contrasena string `json:"contrasena"`
	Telefono   string `json:"telefono"`
}

type UsuariosJwt struct {
	IdUsuario   int       `json:"id_usuario"`
	Token       string    `json:"token"`
	FechaInicio time.Time `json:"fecha_inicio"`
}

type UsuariosOtp struct {
	IdUsuario     int       `json:"id_usuario"`
	CodigoOtp     string    `json:"codigo_otp"`
	FechaCreacion time.Time `json:"fecha_creacion"`
}

// Métodos
func (UsuariosJwt) TableName() string {
	return "usuarios_jwt"
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
