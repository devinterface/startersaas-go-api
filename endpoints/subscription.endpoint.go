package endpoints

import (
	"encoding/json"
	"io/ioutil"

	"devinterface.com/startersaas-go-api/models"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

// SubscriptionEndpoint struct
type SubscriptionEndpoint struct{ BaseEndpoint }

// Subscribe function
func (subscriptionEndpoint *SubscriptionEndpoint) Subscribe(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); can != true {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)
	_, err := govalidator.ValidateMap(inputMap, map[string]interface{}{
		"sourceToken": "ascii,required",
		"planId":      "ascii,required",
	})
	if err != nil {
		return ctx.Status(422).JSON(err.Error())
	}

	me, _ := userEndpoint.CurrentUser(ctx)
	subscription, err := subscriptionService.Subscribe(me.ID, inputMap["planId"].(string), inputMap["sourceToken"].(string))
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(subscription)
}

// GetCustomer function
func (subscriptionEndpoint *SubscriptionEndpoint) GetCustomer(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); can != true {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	account, _ := userEndpoint.CurrentAccount(ctx)
	sCustomer, err := subscriptionService.GetCustomer(account.ID)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(sCustomer)
}

// GetCustomerInvoices function
func (subscriptionEndpoint *SubscriptionEndpoint) GetCustomerInvoices(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); can != true {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	account, _ := userEndpoint.CurrentAccount(ctx)
	sCustomerInvoices, err := subscriptionService.GetCustomerInvoices(account.ID)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(sCustomerInvoices)
}

// GetCustomerCards function
func (subscriptionEndpoint *SubscriptionEndpoint) GetCustomerCards(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); can != true {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	account, _ := userEndpoint.CurrentAccount(ctx)
	sCustomerCards, err := subscriptionService.GetCustomerCards(account.ID)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(sCustomerCards)
}

// CancelSubscription function
func (subscriptionEndpoint *SubscriptionEndpoint) CancelSubscription(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)
	_, err := govalidator.ValidateMap(inputMap, map[string]interface{}{
		"subscriptionId": "ascii,required",
	})
	account, _ := userEndpoint.CurrentAccount(ctx)
	sCustomer, err := subscriptionService.CancelSubscription(account.ID, inputMap["subscriptionId"].(string))
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(sCustomer)
}

// AddCreditCard function
func (subscriptionEndpoint *SubscriptionEndpoint) AddCreditCard(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}

	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)
	_, err := govalidator.ValidateMap(inputMap, map[string]interface{}{
		"sourceToken": "ascii,required",
	})
	if err != nil {
		return ctx.Status(422).JSON(err.Error())
	}

	account, _ := userEndpoint.CurrentAccount(ctx)
	sCustomer, err := subscriptionService.AddCreditCard(account.ID, inputMap["sourceToken"].(string))
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(sCustomer)
}

// RemoveCreditCard function
func (subscriptionEndpoint *SubscriptionEndpoint) RemoveCreditCard(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)
	_, err := govalidator.ValidateMap(inputMap, map[string]interface{}{
		"cardId": "ascii,required",
	})
	if err != nil {
		return ctx.Status(422).JSON(err.Error())
	}
	account, _ := userEndpoint.CurrentAccount(ctx)
	sCustomer, err := subscriptionService.RemoveCreditCard(account.ID, inputMap["cardId"].(string))
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(sCustomer)
}

// SetDefaultCreditCard function
func (subscriptionEndpoint *SubscriptionEndpoint) SetDefaultCreditCard(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)
	_, err := govalidator.ValidateMap(inputMap, map[string]interface{}{
		"cardId": "ascii,required",
	})
	if err != nil {
		return ctx.Status(422).JSON(err.Error())
	}
	account, _ := userEndpoint.CurrentAccount(ctx)
	sCustomer, err := subscriptionService.SetDefaultCreditCard(account.ID, inputMap["cardId"].(string))
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(sCustomer)
}

// Plans function
func (subscriptionEndpoint *SubscriptionEndpoint) Plans(ctx *fiber.Ctx) error {
	data, err := ioutil.ReadFile("./stripe.conf.json")
	var payload interface{}
	json.Unmarshal(data, &payload)
	m := payload.(map[string]interface{})
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(m)
}
