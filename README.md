# prueba-tecnica back-end

### 1. Clonar el repositorio

```bash
git clone https://github.com/Krud3/prueba-tecnica.git
cd prueba-tecnica
```
SSH
```bash
git clone git@github.com:Krud3/prueba-tecnica.git
cd prueba-tecnica
```
### 2. Configurar las variables de entorno

Copia el archivo `.env.example` como `.env` en la raíz del proyecto:

```bash
cp .env.example .env
```

> ⚠️ Si el puerto del front-end (Vite) cambió, actualízalo también en el `.env`.

---

## 🐳 Levantar servicios con Docker

Ejecuta los contenedores de Redis y PostgreSQL:

```bash
docker-compose up -d
```

---

## 🗄️ Preparar la base de datos

Ejecuta la migración de la base de datos:

```bash
migrate -path migrations/ -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up
```

> Reemplaza las variables `DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT` y `DB_NAME` con los valores definidos en tu archivo `.env`.

---

## 🧩 Preparar el entorno de Go (si es necesario)

```bash
go mod tidy
```

---

## 📚 Generar documentación con Swagger (si es necesario)

```bash
swag init -g cmd/api/main.go
```

---

## ▶️ Ejecutar la aplicación

```bash
go run cmd/api/main.go
```

---

## 🏗️ Alternativa con Makefile

Si estás en un entorno compatible con Makefile (como Bash), después de copiar `.env.example` a `.env` puedes ejecutar:

```bash
make up
make migrate
make run-app
```

En el navegador estará swagger
```bash
http://localhost:3000/api/v1/swagger/index.html
```

```bash
prueba-tecnica
├─ Makefile
├─ README.md
├─ cmd
│  └─ api
│     └─ main.go
├─ docker-compose.yml
├─ docs
│  ├─ docs.go
│  ├─ swagger.json
│  └─ swagger.yaml
├─ go.mod
├─ go.sum
├─ internal
│  ├─ adapters
│  │  ├─ rest
│  │  │  ├─ customer_handler.go
│  │  │  ├─ dto.go
│  │  │  ├─ router.go
│  │  │  └─ workorder_handler.go
│  │  └─ storage
│  │     ├─ customer_repository.go
│  │     ├─ db.go
│  │     └─ workorder_repository.go
│  └─ core
│     ├─ domain
│     │  ├─ customer.go
│     │  └─ workorder.go
│     ├─ ports
│     │  └─ ports.go
│     └─ services
│        └─ services.go
└─ migrations
   ├─ 001_create_initial_tables.down.sql
   └─ 001_create_initial_tables.up.sql

```