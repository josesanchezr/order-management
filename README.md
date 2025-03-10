# Order Management Service

Este proyecto es un servicio backend para la gestión de órdenes, desarrollado en **Golang** utilizando **GORM** como ORM y soportando **MySQL** y **Redis**.

## 📌 Requisitos

Antes de ejecutar el proyecto, asegúrate de tener instalados los siguientes programas:

- **Go** (versión 1.19 o superior)
- **MySQL** (versión 8.0 o superior)
- **Redis** (versión 6.0 o superior)
- **Docker** (opcional, para ejecutar MySQL y Redis con contenedores)

---

## 🔧 Configuración de la Base de Datos MySQL

### **1️⃣ Opción 1: Usando Docker**

Si prefieres ejecutar MySQL en un contenedor Docker, usa el siguiente comando:

```sh
docker run --name mysql-container -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=order_db -p 3306:3306 -d mysql:8.0
```

### **2️⃣ Opción 2: Instalación Local**

Si MySQL está instalado localmente, crea la base de datos manualmente:

```sql
CREATE DATABASE order_db;
CREATE USER 'order_user'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON order_db.* TO 'order_user'@'%';
FLUSH PRIVILEGES;
```

Asegúrate de actualizar la configuración en el archivo `.env`:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=order_user
DB_PASSWORD=password
DB_NAME=order_db
```

---

## ⚡ Configuración de Redis

Redis se usa para almacenamiento en caché. Puedes ejecutarlo con:

### **1️⃣ Opción 1: Usando Docker**

```sh
docker run --name redis-container -p 6379:6379 -d redis
```

### **2️⃣ Opción 2: Instalación Local**

Si tienes Redis instalado localmente, simplemente inicia el servicio:

```sh
redis-server
```

Asegúrate de que tu archivo `.env` contenga la configuración correcta:

```env
REDIS_HOST=localhost
REDIS_PORT=6379
```

---

## 🚀 Cómo Ejecutar el Proyecto

1️⃣ **Clonar el repositorio:**

```sh
git clone https://github.com/tuusuario/order-management.git
cd order-management
```

2️⃣ **Instalar dependencias:**

```sh
go mod tidy
```

3️⃣ **Iniciar el servidor:**

```sh
go run cmd/main.go
```

Por defecto, el servicio se ejecutará en `http://localhost:8080`.

---

## 🚀 Iniciar la Aplicación en Local con Docker Compose

1️⃣ **Para iniciar la aplicación con Docker Compose, ejecuta:**

```sh
docker-compose up --build
```

La aplicación estará disponible en `http://localhost:8080`.

---

## 📌 Endpoints Principales

| Método | Endpoint                  | Descripción               |
| ------ | ------------------------- | ------------------------- |
| GET    | `/api/products`           | Lista todas productos     |
| PUT    | `/api/products/:id/stock` | Lista todas las órdenes   |
| POST   | `/api/orders`             | Crea una nueva orden      |
| GET    | `/api/orders/:id`         | Obtiene detalles de orden |

---

## 🧪 Ejecutar Pruebas

Para correr los tests unitarios:

```sh
go test ./... -v
```

Para generar informes de cobertura en **HTML**:

```sh
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

Para generar informes de cobertura en **XML** (para herramientas como **SonarQube**):

```sh
go test -coverprofile=coverage.out ./...
gocov convert coverage.out | gocov-xml > coverage.xml
```

---

## 📄 Licencia

Este proyecto está bajo la licencia **MIT**.
