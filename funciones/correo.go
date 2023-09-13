package funciones

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

type Correo struct {
	Origen  string
	Destino string
	Asunto  string
	Cuerpo  bytes.Buffer
	Mime    string
}

var (
	CORREO    = GetDotEnvVar("CORREO")
	COTRASENA = GetDotEnvVar("CONTRASENA_CORREO")
	SERVIDOR  = GetDotEnvVar("SERVIDOR_SMPT")
	PUERTO    = GetDotEnvVar("PUERTO_CORREO")

	ASUNTO = "Código de verificación"
)

func EnviaCorreoOTPContrasena(destino string, otp string) error {
	var credenciales smtp.Auth = AuthCorreo() //Se obtienen las credenciales para enviar el correo
	var cuerpo bytes.Buffer = bytes.Buffer{}
	var peticion Correo = Correo{}
	var contenido string

	ruta, err := os.Getwd() //Se obtiene la ruta de la carpeta del proyecto para obtener la plantilla
	if err != nil {
		return fmt.Errorf("error al obtener la ruta actual: %v", err)
	}

	plantillaOtp, err := template.ParseFiles(ruta + "/plantillas/otp.html") //Se obtiene la plantilla
	if err != nil {
		return fmt.Errorf("error al obtener la plantilla otp: %v", err)
	}

	plantillaOtp.Execute(&cuerpo, struct { //Se incrusta la información a la plantilla
		Otp string
	}{
		Otp: otp,
	})

	peticion = Correo{
		Origen:  CORREO,
		Destino: destino,
		Asunto:  ASUNTO,
		Cuerpo:  cuerpo,
		Mime:    "1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n",
	}

	contenido = ConstruyeCorreo(peticion)

	/*
		cuerpo := []byte("From: " + "soporte@predicto.ddns.net" + "\r\n" +
			"To: " + destino + "\r\n" +
			"Subject: " + ASUNTO + "\r\n\r\n" +
			"" + otp + "\r\n")*/

	//err = smtp.SendMail(SERVIDOR+":"+PUERTO, credenciales, CORREO, destinos, cuerpo.Bytes()) //El correo es enviado
	err = smtp.SendMail(fmt.Sprintf("%v:%v", SERVIDOR, PUERTO), credenciales, CORREO, []string{peticion.Destino}, []byte(contenido))
	if err != nil {
		return fmt.Errorf("error al enviar correo para el código otp: %v", err)
	}

	return nil
}

// Función para obtener las credenciales del servidor SMTP
func AuthCorreo() smtp.Auth {
	return smtp.PlainAuth("", CORREO, COTRASENA, SERVIDOR)
}

// Función para construir la estructura de un correo
func ConstruyeCorreo(correo Correo) string {
	msg := ""
	msg += fmt.Sprintf("From: %s\r\n", correo.Origen)
	msg += fmt.Sprintf("To: %s\r\n", correo.Destino)
	msg += fmt.Sprintf("Subject: %s\r\n", correo.Asunto)
	msg += fmt.Sprintf("MIME-version: %s\r\n", correo.Mime)
	msg += fmt.Sprintf("\r\n%s\r\n", correo.Cuerpo.String())

	return msg
}
