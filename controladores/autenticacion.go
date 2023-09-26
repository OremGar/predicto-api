package controladores

import (
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/OremGar/predicto-api/bd"
	"github.com/OremGar/predicto-api/configuraciones"
	"github.com/OremGar/predicto-api/funciones"
	"github.com/OremGar/predicto-api/modelos"
	"github.com/OremGar/predicto-api/respuestas"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Función para dar de alta usuarios
func SignUp(w http.ResponseWriter, r *http.Request) {
	var usuario modelos.Usuarios = modelos.Usuarios{}
	var existe bool

	var db *gorm.DB = bd.ConnectDB()
	sqldb, _ := db.DB()
	defer sqldb.Close()

	usuario.Nombre = r.FormValue("nombre")
	usuario.Apellidos = r.FormValue("apellidos")
	usuario.Correo = r.FormValue("correo")
	_, err := mail.ParseAddress(usuario.Correo)
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("el correo no está en el formato correcto"))
		return
	}
	usuario.Usuario = r.FormValue("usuario")
	contrasena, err := funciones.HashContrasena(r.FormValue("contrasena"))
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, err)
		return
	}
	usuario.Contrasena = contrasena
	usuario.Telefono = r.FormValue("telefono")

	err = modelos.ValidarInfoUsuarios(&usuario)
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 100, err)
		return
	}

	result := db.Model(&usuario).Select("count(*) > 0").Where("usuario = ?", usuario.Usuario).Find(&existe)
	if result.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, result.Error)
		return
	}
	if existe {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("el usuario '%v' ya existe", usuario.Usuario))
		return
	}

	result = db.Model(&usuario).Select("count(*) > 0").Where("correo = ?", usuario.Correo).Find(&existe)
	if result.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, result.Error)
		return
	}
	if existe {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("el correo '%v' ya existe", usuario.Correo))
		return
	}

	result = db.Save(&usuario).First(&usuario)
	if result.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, result.Error)
		return
	}

	respuestas.JsonResponse(w, http.StatusCreated, usuario.Id, 0, nil)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var usuario modelos.Usuarios = modelos.Usuarios{}
	var contrasenaPeticion = r.FormValue("contrasena")
	var registroJWT modelos.UsuariosJwt = modelos.UsuariosJwt{}

	usuario.Usuario = r.FormValue("usuario")
	usuario.Correo = r.FormValue("correo")

	var db *gorm.DB = bd.ConnectDB()
	sqldb, _ := db.DB()
	defer sqldb.Close()

	result := db.Model(&usuario).Where("usuario = ? OR correo = ?", usuario.Usuario, usuario.Correo).First(&usuario)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			respuestas.SetError(w, http.StatusNotFound, 100, fmt.Errorf("el usuario no existe"))
			return
		}
		respuestas.SetError(w, http.StatusInternalServerError, 100, result.Error)
		return
	}

	validado := funciones.ValidaContrasena(contrasenaPeticion, usuario.Contrasena)
	if !validado {
		respuestas.SetError(w, http.StatusNotFound, 100, fmt.Errorf("contraseña incorrecta"))
		return
	}

	token, err := configuraciones.GenerarJWT(usuario.Id)
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, err)
		return
	}

	registroJWT = modelos.UsuariosJwt{
		IdUsuario:   usuario.Id,
		Token:       token,
		FechaInicio: time.Now(),
	}

	result = db.Create(&registroJWT)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	r.Header.Set("Authentication", token)

	fmt.Println(token)

	respuestas.JsonResponse(w, http.StatusOK, registroJWT, 0, nil)
}

func RecuperaContrasena(w http.ResponseWriter, r *http.Request) {
	var correo string = r.FormValue("correo")
	var usuario modelos.Usuarios = modelos.Usuarios{}
	var usuarioOtp modelos.UsuariosOtp = modelos.UsuariosOtp{}
	var otp string

	var db *gorm.DB = bd.ConnectDB()
	sqldb, _ := db.DB()
	defer sqldb.Close()

	if correo == "" {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("la solicitud no incluye correo"))
		return
	}

	resultado := db.Model(usuario).Where("correo = ?", correo).First(&usuario)
	if resultado.Error != nil {
		if resultado.Error == gorm.ErrRecordNotFound {
			respuestas.SetError(w, http.StatusNotFound, 100, fmt.Errorf("el correo ingresado no existe en el sistema"))
			return
		}
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("problema interno: %v", resultado.Error))
		return
	}

	otp = funciones.GeneraOTP(12)
	err := funciones.EnviaCorreoOTPContrasena(correo, otp)
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error al enviar correo: %v", err))
		return
	}

	usuarioOtp = modelos.UsuariosOtp{
		IdUsuario:     usuario.Id,
		CodigoOtp:     otp,
		FechaCreacion: time.Now(),
	}

	resultado = db.Create(&usuarioOtp)
	if resultado.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error al guardar otp: %v", err))
		return
	}

	respuestas.JsonResponse(w, http.StatusOK, nil, 0, nil)
}

func NuevaContrasena(w http.ResponseWriter, r *http.Request) {
	var nuevaContrasena string = r.FormValue("nuevaContrasena")
	var otp string = r.FormValue("otp")

	var db *gorm.DB = bd.ConnectDB()
	sqldb, _ := db.DB()
	defer sqldb.Close()

	_, otp_usuario, err := modelos.ChecarSiOTPValido(otp)
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 100, err)
		return
	}

	usuario, err := modelos.ChecarSiUsuarioExiste(otp_usuario.IdUsuario)
	if err != nil {
		respuestas.SetError(w, http.StatusBadRequest, 100, err)
		return
	}

	usuario.Contrasena, err = funciones.HashContrasena(nuevaContrasena)
	if err != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, err)
		return
	}

	result := db.Save(&usuario)
	if result.Error != nil {
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error interno sql: %v", result.Error))
		return
	}

	respuestas.JsonResponse(w, http.StatusOK, usuario, 0, nil)
}

func ValidaOTPNvaContrasena(w http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var codigoOTP string = vars["codigo"]
	var usuarioOTP modelos.UsuariosOtp = modelos.UsuariosOtp{}

	var db *gorm.DB = bd.ConnectDB()
	sqldb, _ := db.DB()
	defer sqldb.Close()

	if codigoOTP == "" {
		respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("no se incluyó el código OTP en la url"))
		return
	}

	result := db.Model(&usuarioOTP).Where("codigo_otp = ?", codigoOTP).First(&usuarioOTP)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			respuestas.SetError(w, http.StatusBadRequest, 100, fmt.Errorf("no existe código"))
			return
		}
		respuestas.SetError(w, http.StatusInternalServerError, 100, fmt.Errorf("error en la consulta: %v", result.Error))
		return
	}

	respuestas.JsonResponse(w, http.StatusOK, 100, 0, nil)
}

func Prueba(w http.ResponseWriter, _ *http.Request) {
	respuestas.JsonResponse(w, http.StatusOK, "Saludos", 0, nil)
}
