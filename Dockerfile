# Etapa de compilación
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copiar y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# 🚀 Compilar de forma estática (sin dependencias del sistema)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

RUN ls -l /app/main

# Etapa de ejecución
FROM alpine:latest

WORKDIR /app

# Instalar dependencias necesarias para Alpine
RUN apk add --no-cache ca-certificates

# Copiar el binario desde la etapa de compilación
COPY --from=builder /app/main .

# Verificar el binario copiado
RUN ls -l /app/main

# Asegurar permisos de ejecución
RUN chmod +x /app/main

# Ejecutar la aplicación
CMD ["/app/main"]
