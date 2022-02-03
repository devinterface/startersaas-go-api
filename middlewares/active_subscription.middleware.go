package middlewares

import (
	"devinterface.com/startersaas-go-api/models"
	"github.com/gofiber/fiber/v2"
)

var ActiveSubscription = func(ctx *fiber.Ctx) (err error) {
	currentAccount := ctx.Locals("currentAccount").(*models.Account)
	if currentAccount.SubscriptionStatus() != models.SubscriptionDeactivated {
		return ctx.Next()
	} else {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
}
