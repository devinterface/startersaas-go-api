package main

import (
	"log"
	"os"
	"time"

	"devinterface.com/startersaas-go-api/endpoints"
	"devinterface.com/startersaas-go-api/services"
	"github.com/Kamva/mgm/v3"
	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initDabatase() {
	mgm.SetDefaultConfig(nil, os.Getenv("DATABASE"), options.Client().ApplyURI(os.Getenv("DATABASE_URI")))
}

func runScheduledNotifications() (err error) {
	var subscriptionService = services.SubscriptionService{}
	subscriptionService.RunNotifyExpiringTrials()
	subscriptionService.RunNotifyPaymentFailed()
	return err
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

	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Day().At("00:01").Do(runScheduledNotifications)
	s.StartAsync()

	log.Fatal(app.Listen(os.Getenv("PORT")))
}
