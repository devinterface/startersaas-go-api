package endpoints

import (
	"github.com/gofiber/fiber/v2"
)

var authEndpoint = AuthEndpoint{}
var accountEndpoint = AccountEndpoint{}
var userEndpoint = UserEndpoint{}
var subscriptionEndpoint = SubscriptionEndpoint{}
var webhookEndpoint = WebhookEndpoint{}

// SetupPublicRoutes function
func SetupPublicRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/auth/login", authEndpoint.Login)
	v1.Post("/auth/signup", authEndpoint.Signup)
	v1.Post("/auth/send-activation-link", authEndpoint.SendActivationLink)
	v1.Post("/auth/activate", authEndpoint.Activate)
	v1.Post("/auth/send-forgot-password-link", authEndpoint.SendForgotPasswordLink)
	v1.Post("/auth/reset-password", authEndpoint.ResetPassword)
	v1.Post("/auth/sso-login", authEndpoint.SsoLogin)
	v1.Post("/stripe/webhook", webhookEndpoint.HandleWebhook)
	v1.Post("/fattura24", webhookEndpoint.Fattura24)
}

// SetupPrivateRoutes function
func SetupPrivateRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/auth/refresh-token", authEndpoint.RefreshToken)
	v1.Get("/accounts/:id", accountEndpoint.ByID)
	v1.Put("/accounts/:id", accountEndpoint.Update)
	v1.Get("/users/me", userEndpoint.Me)
	v1.Put("/users/me", userEndpoint.UpdateMe)
	v1.Put("/users/me/change-password", userEndpoint.ChangePassword)
	v1.Put("/users/me/generate-sso", userEndpoint.GenerateSso)
	v1.Get("/users", userEndpoint.Index)
	v1.Post("/users", userEndpoint.Create)
	v1.Put("/users/:id", userEndpoint.Update)
	v1.Delete("/users/:id", userEndpoint.Delete)
	v1.Post("/stripe/subscriptions", subscriptionEndpoint.Subscribe)
	v1.Delete("/stripe/subscriptions", subscriptionEndpoint.CancelSubscription)
	v1.Get("/stripe/customers/me", subscriptionEndpoint.GetCustomer)
	v1.Get("/stripe/customers/me/invoices", subscriptionEndpoint.GetCustomerInvoices)
	v1.Get("/stripe/customers/me/cards", subscriptionEndpoint.GetCustomerCards)
	v1.Post("/stripe/cards", subscriptionEndpoint.AddCreditCard)
	v1.Delete("/stripe/cards", subscriptionEndpoint.RemoveCreditCard)
	v1.Put("/stripe/cards", subscriptionEndpoint.SetDefaultCreditCard)
}
