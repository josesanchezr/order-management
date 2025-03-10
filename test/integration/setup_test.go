package integration_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"order_management/internal/models"
	"order_management/internal/validators"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB
	server         *httptest.Server
	mysqlContainer testcontainers.Container
	redisClient    *redis.Client
)

// SetupTestDatabase inicializa un contenedor MySQL y devuelve una conexión GORM.
func SetupTestDatabase(t *testing.T) *gorm.DB {
	ctx := context.Background()

	// Configuración del contenedor MySQL
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      "testdb",
			"MYSQL_USER":          "testuser",
			"MYSQL_PASSWORD":      "testpass",
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server").WithStartupTimeout(30 * time.Second),
	}

	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Error iniciando el contenedor de MySQL: %v", err)
	}
	mysqlContainer = mysqlC

	// Obtener host y puerto de MySQL
	host, err := mysqlC.Host(ctx)
	if err != nil {
		t.Fatalf("Error obteniendo el host del contenedor: %v", err)
	}
	port, err := mysqlC.MappedPort(ctx, "3306")
	if err != nil {
		t.Fatalf("Error obteniendo el puerto del contenedor: %v", err)
	}

	// Crear conexión con MySQL
	dsn := fmt.Sprintf("testuser:testpass@tcp(%s:%s)/testdb?charset=utf8mb4&parseTime=True&loc=Local", host, port.Port())
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error conectando a MySQL: %v", err)
	}

	// Migrar modelos
	err = db.AutoMigrate(&models.Product{}, &models.Order{}, &models.OrderItem{})
	if err != nil {
		t.Fatalf("Error ejecutando migraciones: %v", err)
	}

	return db
}

// SetupTestRedis inicializa un cliente Redis.
func SetupTestRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // También se puede usar Testcontainers para Redis si se requiere
		DB:   0,
	})
	return client
}

// SetupTestServer configura el servidor de pruebas con Echo y GORM.
func SetupTestServer(t *testing.T, registerRoutes func(e *echo.Echo, db *gorm.DB, redisClient *redis.Client)) {
	db = SetupTestDatabase(t)
	redisClient = SetupTestRedis()

	// Crear Echo y registrar rutas
	e := echo.New()

	// Configurar el validador
	e.Validator = validators.NewValidator()

	registerRoutes(e, db, redisClient)

	// Iniciar servidor de pruebas
	server = httptest.NewServer(e)
}

// TearDown detiene el servidor y el contenedor de MySQL después de cada test.
func TearDown() {
	if server != nil {
		server.Close()
	}
	if mysqlContainer != nil {
		mysqlContainer.Terminate(context.Background())
	}
	if redisClient != nil {
		redisClient.Close()
	}
}
