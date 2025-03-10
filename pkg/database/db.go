package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB inicializa y devuelve una conexión a la base de datos MySQL
func InitDB() *gorm.DB {
	// Obtener las variables de entorno para la conexión a MySQL
	dbUser := getEnv("DB_USER", "order")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "order_management")

	// Construir la cadena de conexión
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	// Conectar a MySQL usando GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	log.Println("Conexión a la base de datos MySQL establecida exitosamente")
	return db
}

// getEnv obtiene una variable de entorno o un valor por defecto si no está definida
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
