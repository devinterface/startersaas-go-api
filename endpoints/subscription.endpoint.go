package endpoints

import (
	"encoding/json"
	"io/ioutil"

	"devinterface.com/startersaas-go-api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

// SubscriptionEndpoint struct
type SubscriptionEndpoint struct{ BaseEndpoint }

// Subscribe function
func (subscriptionEndpoint *SubscriptionEndpoint) Subscribe(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)

	v := validate.Map(inputMap)

	v.StringRule("planId", "ascii|required")

	if !v.Validate() {
		return ctx.Status(422).JSON(v.Errors)
	}

	me, _ := userEndpoint.CurrentUser(ctx)
	subscription, err := subscriptionService.Subscribe(me.ID, inputMap["planId"].(string))
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	marshalled, _ := json.Marshal(subscription)
	var payload interface{}
	json.Unmarshal(marshalled, &payload)
	return ctx.JSON(payload)
}

// GetCustomer function
func (subscriptionEndpoint *SubscriptionEndpoint) GetCustomer(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
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
	marshalled, _ := json.Marshal(sCustomer)
	var payload interface{}
	json.Unmarshal(marshalled, &payload)
	return ctx.Status(200).JSON(payload)
}

// GetCustomerInvoices function
func (subscriptionEndpoint *SubscriptionEndpoint) GetCustomerInvoices(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	account, _ := userEndpoint.CurrentAccount(ctx)
	sCustomerInvoices, err := subscriptionService.GetCustomerInvoices(account.ID)
	if err != nil {
		ctx.JSON([]string{})
	}
	if len(sCustomerInvoices) == 0 {
		return ctx.JSON([]string{})
	}
	marshalled, _ := json.Marshal(sCustomerInvoices)
	var payload interface{}
	json.Unmarshal(marshalled, &payload)
	return ctx.Status(200).JSON(payload)
}

// GetCustomerCards function
func (subscriptionEndpoint *SubscriptionEndpoint) GetCustomerCards(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	account, _ := userEndpoint.CurrentAccount(ctx)
	sCustomerCards, err := subscriptionService.GetCustomerCards(account.ID)
	if err != nil {
		ctx.JSON([]string{})
	}
	if len(sCustomerCards) == 0 {
		return ctx.JSON([]string{})
	}
	marshalled, _ := json.Marshal(sCustomerCards)
	var payload interface{}
	json.Unmarshal(marshalled, &payload)
	return ctx.Status(200).JSON(payload)
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
	v := validate.Map(inputMap)
	v.StringRule("subscriptionId", "ascii|required")

	if !v.Validate() {
		return ctx.Status(422).JSON(v.Errors)
	}
	account, _ := userEndpoint.CurrentAccount(ctx)
	sCustomer, err := subscriptionService.CancelSubscription(account.ID, inputMap["subscriptionId"].(string))
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	marshalled, _ := json.Marshal(sCustomer)
	var payload interface{}
	json.Unmarshal(marshalled, &payload)
	return ctx.Status(200).JSON(payload)
}

// AddCreditCard function
func (subscriptionEndpoint *SubscriptionEndpoint) CreateSetupIntent(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}

	account, _ := userEndpoint.CurrentAccount(ctx)
	setupIntent, err := subscriptionService.CreateSetupIntent(account.ID)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	marshalled, _ := json.Marshal(setupIntent)
	var payload interface{}
	json.Unmarshal(marshalled, &payload)
	return ctx.Status(200).JSON(payload)
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
	v := validate.Map(inputMap)
	v.StringRule("cardId", "ascii|required")

	if !v.Validate() {
		return ctx.Status(422).JSON(v.Errors)
	}
	account, _ := userEndpoint.CurrentAccount(ctx)
	sCustomer, err := subscriptionService.RemoveCreditCard(account.ID, inputMap["cardId"].(string))
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	marshalled, _ := json.Marshal(sCustomer)
	var payload interface{}
	json.Unmarshal(marshalled, &payload)
	return ctx.Status(200).JSON(payload)
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
	v := validate.Map(inputMap)
	v.StringRule("cardId", "ascii|required")

	if !v.Validate() {
		return ctx.Status(422).JSON(v.Errors)
	}
	account, _ := userEndpoint.CurrentAccount(ctx)
	sCustomer, err := subscriptionService.SetDefaultCreditCard(account.ID, inputMap["cardId"].(string))
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	marshalled, _ := json.Marshal(sCustomer)
	var payload interface{}
	json.Unmarshal(marshalled, &payload)
	return ctx.Status(200).JSON(payload)
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
	return ctx.Status(200).JSON(m)
}

// CreateCustomerCheckoutSession function
func (subscriptionEndpoint *SubscriptionEndpoint) CreateCustomerCheckoutSession(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)
	v := validate.Map(inputMap)
	v.StringRule("planId", "ascii|required")

	if !v.Validate() {
		return ctx.Status(422).JSON(v.Errors)
	}
	me, _ := userEndpoint.CurrentUser(ctx)
	redirectUrl, err := subscriptionService.CreateCustomerCheckoutSession(me.ID, inputMap["planId"].(string))
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"redirect_url": redirectUrl,
	})
}

// CreateCustomerPortalSession function
func (subscriptionEndpoint *SubscriptionEndpoint) CreateCustomerPortalSession(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	me, _ := userEndpoint.CurrentUser(ctx)
	redirectUrl, err := subscriptionService.CreateCustomerPortalSession(me.AccountID)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"redirect_url": redirectUrl,
	})
}
