package middlewares

import (
	"strings"

	"devinterface.com/startersaas-go-api/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var accountService = services.AccountService{}

// APITokenAuth middleware
var APITokenAuth = func(ctx *fiber.Ctx) (err error) {
	bearer := ctx.Get("authorization")
	if len(bearer) > 0 {
		components := strings.SplitN(bearer, " ", 2)
		if len(components) == 2 {
			apiToken := components[1]
			if apiToken != "" {
				account, _ := accountService.OneBy(bson.M{"apiToken": apiToken})
				if account.ID != primitive.NilObjectID {
					ctx.Locals("currentAccount", account)
					return ctx.Next()
				}
			}
		}
	}
	return ctx.Status(401).JSON(fiber.Map{
		"message": "You are not authorized to perform this action",
	})
}
