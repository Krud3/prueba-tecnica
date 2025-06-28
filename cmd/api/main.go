// cmd/api/main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/krud3/prueba-tecnica/internal/adapters/storage"
)

func main() {
	db, err := storage.NewGormDB()
	if err != nil {
		log.Fatalf("Error al conectar a la DB: %v", err)
	}

	app := fiber.New()

	// testing ping
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	// testing db
	app.Get("/db", func(c *fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err != nil {
			return c.Status(500).SendString("Error obteniendo DB subestructura")
		}

		if err := sqlDB.Ping(); err != nil {
			return c.Status(500).SendString("No se pudo conectar a la base de datos")
		}

		return c.SendString("Base de datos conectada correctamente")
	})

	log.Println("Servidor escuchando en :3000")
	log.Fatal(app.Listen(":3000"))

}
