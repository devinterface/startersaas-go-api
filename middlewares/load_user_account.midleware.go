package middlewares

import (
	"devinterface.com/startersaas-go-api/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

var userService = services.UserService{}

var LoadUserAccount = func(ctx *fiber.Ctx) (err error) {
	jwtUser := ctx.Locals("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	q := bson.M{"email": claims["email"].(string)}
	me, err := userService.OneBy(q)
	if err != nil {
		return err
	}
	account, err := accountService.ByID(me.AccountID)
	if err != nil {
		return err
	}
	ctx.Locals("currentUser", me)
	ctx.Locals("currentAccount", account)
	return ctx.Next()

}
