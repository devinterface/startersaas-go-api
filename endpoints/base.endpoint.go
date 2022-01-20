package endpoints

import (
	"devinterface.com/startersaas-go-api/models"
	"devinterface.com/startersaas-go-api/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

// BaseEndpoint struct
type BaseEndpoint struct{}

var authService = services.AuthService{}
var accountService = services.AccountService{}
var userService = services.UserService{}
var emailService = services.EmailService{}
var subscriptionService = services.SubscriptionService{}
var webhookService = services.WebhookService{}
var fattura24Service = services.Fattura24Service{}

// CurrentUser function
func (baseEndpoint *BaseEndpoint) CurrentUser(ctx *fiber.Ctx) (me *models.User, err error) {
	jwtUser := ctx.Locals("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	q := bson.M{"email": claims["email"].(string)}
	me, err = userService.OneBy(q)
	return me, err
}

// CurrentAccount function
func (baseEndpoint *BaseEndpoint) CurrentAccount(ctx *fiber.Ctx) (currentAccount *models.Account, err error) {
	currentUser, err := baseEndpoint.CurrentUser(ctx)
	q := bson.M{"_id": currentUser.AccountID}
	currentAccount, err = accountService.OneBy(q)
	return currentAccount, err
}

// Can function
func (baseEndpoint *BaseEndpoint) Can(ctx *fiber.Ctx, role string) (success bool) {
	jwtUser := ctx.Locals("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	q := bson.M{"email": claims["email"].(string)}
	me, _ := userService.OneBy(q)
	return me.Role == role
}

func (baseEndpoint *BaseEndpoint) ActiveSubscription(ctx *fiber.Ctx) (isActive bool) {
	currentAccount, _ := baseEndpoint.CurrentAccount(ctx)
	return currentAccount.SubscriptionStatus() != models.SubscriptionDeactivated || currentAccount.SubscriptionStatus() != models.SubscriptionPendingActivation
}

func buildMeta(page int64, limit int64, count int64) (meta map[string]int64) {
	meta = make(map[string]int64)
	meta["page"] = page
	meta["limit"] = limit
	meta["count"] = count
	if page > 1 {
		meta["prev"] = page - 1
	}
	if page*limit < count {
		meta["next"] = page + 1
	}
	return meta
}
