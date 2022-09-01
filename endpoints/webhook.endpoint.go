package endpoints

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v72"
	"go.mongodb.org/mongo-driver/bson"
)

// WebhookEndpoint struct
type WebhookEndpoint struct{ BaseEndpoint }

// HandleWebhook function
func (webhookEndpoint *WebhookEndpoint) HandleWebhook(ctx *fiber.Ctx) error {
	event := stripe.Event{}
	ctx.BodyParser(&event)
	payload := make(map[string]interface{})
	ctx.BodyParser(&payload)
	go webhookService.HandleWebhook(payload, event)
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func (webhookEndpoint *WebhookEndpoint) Fattura24(ctx *fiber.Ctx) error {
	account, _ := accountService.OneBy(bson.M{"subdomain": "devinterface"})
	fattura24Service.DummyGenerateInvoice(account.ID)
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}
