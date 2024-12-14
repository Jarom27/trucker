# Etapa 1: Construcción
FROM golang:1.20 AS builder

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar los archivos del proyecto al contenedor
COPY . .

# Descargar las dependencias y compilar la aplicación
RUN go mod tidy
RUN go build -o server .

# Etapa 2: Imagen final
FROM debian:bullseye-slim

# Crear un usuario no root (opcional)
RUN useradd -m tcpuser

# Copiar el binario compilado desde la etapa de construcción
COPY --from=builder /app/server /usr/local/bin/server

# Cambiar al usuario no root
USER tcpuser

# Exponer el puerto en el que escucha tu servidor (en este caso, 7700)
EXPOSE 7700

# Comando para ejecutar la aplicación
CMD ["server"]
