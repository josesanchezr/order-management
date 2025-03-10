package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

const (
	KeyPrefix        = "idempotency:"
	StatusInProgress = "IN_PROGRESS"
	StatusCompleted  = "COMPLETED"
	TTL              = 10 * time.Minute
)

type IdempotencyData struct {
	Status   string      `json:"status"`
	Response interface{} `json:"response"`
}

// IdempotencyMiddleware contiene la lógica de idempotencia al crear una orden
func IdempotencyMiddleware(redisClient *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.Background()
			idempotencyKey := c.Request().Header.Get("Idempotency-Key")

			// Si no tiene Idempotency-Key, continuar sin usar Redis
			if idempotencyKey == "" {
				return next(c)
			}

			redisKey := KeyPrefix + idempotencyKey

			// Revisar si la clave ya existe en Redis
			data, err := redisClient.Get(ctx, redisKey).Result()
			if err == nil {
				var storedData IdempotencyData
				json.Unmarshal([]byte(data), &storedData)

				// Si la solicitud está en progreso, devolver error 409
				if storedData.Status == StatusInProgress {
					return c.JSON(http.StatusConflict, map[string]string{"error": "Petición esta siendo procesada"})
				}

				// Si la solicitud ya fue completada, devolver la respuesta almacenada
				return c.JSON(http.StatusOK, storedData.Response)
			}

			// Guardar el estado IN_PROGRESS en Redis antes de procesar la solicitud
			storedData := IdempotencyData{
				Status:   StatusInProgress,
				Response: nil,
			}
			jsonData, _ := json.Marshal(storedData)
			redisClient.Set(ctx, redisKey, jsonData, TTL)

			// Capturar la respuesta antes de enviarla al cliente
			rec := c.Response().Writer
			buffer := new(bytes.Buffer)
			multiWriter := io.MultiWriter(rec, buffer)
			c.Response().Writer = &responseWriterInterceptor{ResponseWriter: rec, writer: multiWriter}

			// Procesar la solicitud
			err = next(c)

			// Guardar la respuesta en Redis solo si no hay errores
			if err == nil {
				responseData := IdempotencyData{
					Status:   StatusCompleted,
					Response: json.RawMessage(buffer.Bytes()),
				}
				jsonData, _ := json.Marshal(responseData)
				redisClient.Set(ctx, redisKey, jsonData, TTL)
			}

			return err
		}
	}
}

// responseWriterInterceptor captura la respuesta antes de enviarla
type responseWriterInterceptor struct {
	http.ResponseWriter
	writer io.Writer
}

func (r *responseWriterInterceptor) Write(b []byte) (int, error) {
	return r.writer.Write(b)
}
