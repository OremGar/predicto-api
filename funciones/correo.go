package funciones

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

var (
	CORREO    = GetDotEnvVar("CORREO")
	COTRASENA = GetDotEnvVar("CONTRASENA_CORREO")
	SERVIDOR  = GetDotEnvVar("SERVIDOR_SMPT")
	PUERTO    = GetDotEnvVar("PUERTO_CORREO")

	ASUNTO = "Código de verificación"
)

func EnviaCorreoOTP(destino string, otp string) error {
	var credenciales smtp.Auth = AuthCorreo()                                                      //Se obtienen las credenciales para enviar el correo
	var destinos []string = []string{destino}                                                      //Se agrega al único destino al slice
	var cuerpo bytes.Buffer                                                                        //Se crea objeto para añadir información al buffer
	var mimeHeaders string = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" //Encabezados para la plantilla HTML

	cuerpo.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", ASUNTO, mimeHeaders))) //Se añade el asunto y los encabezados al cuerpo

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

	//err = smtp.SendMail(SERVIDOR+":"+PUERTO, credenciales, CORREO, destinos, cuerpo.Bytes()) //El correo es enviado
	err = smtp.SendMail(fmt.Sprintf("%v:%v", "mail.noip.com", 587), credenciales, "soporte@predicto.ddns.net", destinos, cuerpo.Bytes())
	if err != nil {
		return fmt.Errorf("error al enviar correo para el código otp: %v", err)
	}

	return nil
}

// Función para obtener las credenciales del servidor SMTP
func AuthCorreo() smtp.Auth {
	return smtp.PlainAuth("", CORREO, COTRASENA, SERVIDOR)
}
