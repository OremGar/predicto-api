package funciones

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"
	"strconv"

	gomail "gopkg.in/mail.v2"
)

type Correo struct {
	Origen  string
	Destino string
	Asunto  string
	Cuerpo  bytes.Buffer
	Mime    string
}

var (
	CORREO     = GetDotEnvVar("CORREO")
	CONTRASENA = GetDotEnvVar("CONTRASENA_CORREO")
	SERVIDOR   = GetDotEnvVar("SERVIDOR_SMPT")
	PUERTO     = GetDotEnvVar("PUERTO_CORREO")

	ASUNTO = "PREDICTO"
)

func EnviaCorreoOTPContrasena(destino string, otp string) error {
	var m *gomail.Message = gomail.NewMessage()
	var d *gomail.Dialer = &gomail.Dialer{}
	var cuerpo bytes.Buffer = bytes.Buffer{}
	var puerto int

	ruta, err := os.Getwd() //Se obtiene la ruta de la carpeta del proyecto para obtener la plantilla
	if err != nil {
		return fmt.Errorf("error al obtener la ruta actual: %v", err)
	}

	plantillaOtp, err := template.ParseFiles(ruta + "/plantillas/otpcontrasena.html") //Se obtiene la plantilla
	if err != nil {
		return fmt.Errorf("error al obtener la plantilla otp: %v", err)
	}

	if GetDotEnvVar("PRODUCCION") == "true" {
		otp = fmt.Sprintf("https://predicto.ddns.net/RecuperacionContrasena?codigo=%v", otp)
	} else {
		otp = fmt.Sprintf("http://localhost:3000/RecuperacionContrasena?codigo=%v", otp)
	}

	plantillaOtp.Execute(&cuerpo, struct { //Se incrusta la información a la plantilla
		Otp string
	}{
		Otp: otp,
	})

	// Set E-Mail sender
	m.SetHeader("From", CORREO)

	// Set E-Mail receivers
	m.SetHeader("To", destino)

	// Set E-Mail subject
	m.SetHeader("Subject", "Recuperar contraseña")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", cuerpo.String())

	puerto, err = strconv.Atoi(PUERTO)
	if err != nil {
		return fmt.Errorf("error al enviar correo para el código otp: %v", err)
	}
	d = gomail.NewDialer(SERVIDOR, puerto, CORREO, CONTRASENA)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	/*peticion = Correo{
		Origen:  CORREO,
		Destino: destino,
		Asunto:  ASUNTO,
		Cuerpo:  cuerpo,
		Mime:    "1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n",
	}

	contenido = ConstruyeCorreo(peticion)

	err = smtp.SendMail(fmt.Sprintf("%v:%v", SERVIDOR, PUERTO), credenciales, CORREO, []string{peticion.Destino}, []byte(contenido))
	if err != nil {
		return fmt.Errorf("error al enviar correo para el código otp: %v", err)
	}*/

	return nil
}

func EnviaCorreoOTPLogin(destino string, otp string) error {
	var m *gomail.Message = gomail.NewMessage()
	var d *gomail.Dialer = &gomail.Dialer{}
	var cuerpo bytes.Buffer = bytes.Buffer{}
	var puerto int
	var err error

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

	// Set E-Mail sender
	m.SetHeader("From", CORREO)

	// Set E-Mail receivers
	m.SetHeader("To", destino)

	// Set E-Mail subject
	m.SetHeader("Subject", "Doble auntenticación")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", cuerpo.String())

	puerto, err = strconv.Atoi(PUERTO)
	if err != nil {
		return fmt.Errorf("error al enviar correo para el código otp: %v", err)
	}
	d = gomail.NewDialer(SERVIDOR, puerto, CORREO, CONTRASENA)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error al enviar correo para el código otp: %v", err)
	}

	/*
		peticion = Correo{
			Origen:  CORREO,
			Destino: destino,
			Asunto:  ASUNTO,
			Cuerpo:  cuerpo,
			Mime:    "1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n",
		}

		contenido = ConstruyeCorreo(peticion)

		err = smtp.SendMail(fmt.Sprintf("%v:%v", SERVIDOR, PUERTO), credenciales, CORREO, []string{peticion.Destino}, []byte(contenido))
		if err != nil {
			return fmt.Errorf("error al enviar correo para el código otp: %v", err)
		}*/

	return nil
}

// Función para obtener las credenciales del servidor SMTP
/*
func AuthCorreo() smtp.Auth {
	return smtp.PlainAuth("", CORREO, COTRASENA, SERVIDOR)
}*/

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
