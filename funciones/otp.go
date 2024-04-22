package funciones

import (
	"crypto/rand"
	"io"
)

var table [10]byte = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GeneraOTP(limite int) string {
	var b []byte = make([]byte, limite)
	n, err := io.ReadAtLeast(rand.Reader, b, limite)
	if n != limite {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
