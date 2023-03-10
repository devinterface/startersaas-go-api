package main

import (
	"context"
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
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gookit/validate"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initDabatase() {
	monitor := &event.CommandMonitor{
		Started: func(_ context.Context, e *event.CommandStartedEvent) {
			//fmt.Println(e.Command)
		},
		Succeeded: func(_ context.Context, e *event.CommandSucceededEvent) {
			//fmt.Println(e.Reply)
		},
		Failed: func(_ context.Context, e *event.CommandFailedEvent) {
			//fmt.Println(e.Failure)
		},
	}
	opts := options.Client().SetMonitor(monitor)
	mgm.SetDefaultConfig(nil, os.Getenv("DATABASE"), options.Client().ApplyURI(os.Getenv("DATABASE_URI")), opts)
}
func runScheduledNotifications() (err error) {
	var subscriptionService = services.SubscriptionService{}
	subscriptionService.RunNotifyExpiringTrials()
	subscriptionService.RunNotifyPaymentFailed()
	return err
}

func storeEmails() {
	var emailService = services.EmailService{}
	emailService.StoreEmails()
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CORS_SITES"),
	}))
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{Format: "[${time}] ${locals:requestid} ${status} - ${method} ${path}​\n"}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	initDabatase()
	storeEmails()

	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
	})

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
