package main

import (
	"log"
	"os"

	"devinterface.com/goaas-api-starter/endpoints"
	"github.com/Kamva/mgm/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initDabatase() {
	mgm.SetDefaultConfig(nil, os.Getenv("DATABASE"), options.Client().ApplyURI("mongodb://localhost:27017"))
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CORS_SITES"),
	}))
	app.Use(logger.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	initDabatase()

	endpoints.SetupPublicRoutes(app)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	endpoints.SetupPrivateRoutes(app)

	log.Fatal(app.Listen(os.Getenv("PORT")))
}
