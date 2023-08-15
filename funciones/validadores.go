package funciones

import "net/mail"

// Función para validar el formato de un correo
func ValidaCorreo(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
