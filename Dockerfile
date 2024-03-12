#Imagen a descargar
FROM golang:1.22

#Se actualiza el contenedor
RUN apt-get update

#Directorio de trabajo
WORKDIR /go/src

#Copia el proyecto en el directorio actual para posteriormente instalar las librerías usadas en el proyecto
COPY . .
RUN go mod download

#Compila el main.go
RUN go build -o /main

#Expone el puerto
EXPOSE 8081

#Ejecuta el compilado una vez que el contenedor arranca
CMD [ "/main" ]