package middlewares

import (
	"devinterface.com/startersaas-go-api/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var accountService = services.AccountService{}

// AuthReq middleware
var APITokenAuth = func(ctx *fiber.Ctx) (err error) {
	apiToken := ctx.Query("apiToken")
	if apiToken != "" {
		account, _ := accountService.OneBy(bson.M{"apiToken": apiToken})
		if account.ID != primitive.NilObjectID {
			ctx.Locals("currentAccount", account)
			return ctx.Next()
		}
	}
	return ctx.Status(401).JSON(fiber.Map{
		"message": "You are not authorized to perform this action",
	})
}
