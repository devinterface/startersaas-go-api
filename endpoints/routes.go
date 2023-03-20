package endpoints

import (
	"devinterface.com/startersaas-go-api/middlewares"
	"github.com/gofiber/fiber/v2"
)

var authEndpoint = AuthEndpoint{}
var accountEndpoint = AccountEndpoint{}
var userEndpoint = UserEndpoint{}
var subscriptionEndpoint = SubscriptionEndpoint{}
var webhookEndpoint = WebhookEndpoint{}
var teamEndpoint = TeamEndpoint{}

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
	v1.Get("/stripe/plans", subscriptionEndpoint.Plans)
}

// SetupPrivateRoutes function
func SetupPrivateRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/auth/refresh-token", middlewares.LoadUserAccount, authEndpoint.RefreshToken)
	v1.Get("/accounts/:id", middlewares.LoadUserAccount, accountEndpoint.ByID)
	v1.Put("/accounts/:id", middlewares.LoadUserAccount, accountEndpoint.Update)
	v1.Get("/teams", middlewares.LoadUserAccount, teamEndpoint.Index)
	v1.Get("/teams/:id", middlewares.LoadUserAccount, teamEndpoint.ByID)
	v1.Post("/teams", middlewares.LoadUserAccount, teamEndpoint.Create)
	v1.Delete("/teams/:id", middlewares.LoadUserAccount, teamEndpoint.Delete)
	v1.Put("/teams/:id", middlewares.LoadUserAccount, teamEndpoint.Update)
	v1.Put("/teams/:id/add-user/:userId", middlewares.LoadUserAccount, teamEndpoint.AddUser)
	v1.Put("/teams/:id/remove-user/:userId", middlewares.LoadUserAccount, teamEndpoint.RemoveUser)
	v1.Get("/users/me", middlewares.LoadUserAccount, userEndpoint.Me)
	v1.Put("/users/me", middlewares.LoadUserAccount, userEndpoint.UpdateMe)
	v1.Put("/users/me/change-password", middlewares.LoadUserAccount, userEndpoint.ChangePassword)
	v1.Put("/users/me/generate-sso", middlewares.LoadUserAccount, userEndpoint.GenerateSso)
	v1.Get("/users", middlewares.LoadUserAccount, userEndpoint.Index)
	v1.Get("/users/:id", middlewares.LoadUserAccount, userEndpoint.ByID)
	v1.Post("/users", middlewares.LoadUserAccount, userEndpoint.Create)
	v1.Put("/users/:id", middlewares.LoadUserAccount, userEndpoint.Update)
	v1.Delete("/users/:id", middlewares.LoadUserAccount, userEndpoint.Delete)
	v1.Post("/stripe/subscriptions", middlewares.LoadUserAccount, subscriptionEndpoint.Subscribe)
	v1.Delete("/stripe/subscriptions", middlewares.LoadUserAccount, subscriptionEndpoint.CancelSubscription)
	v1.Get("/stripe/customers/me", middlewares.LoadUserAccount, subscriptionEndpoint.GetCustomer)
	v1.Get("/stripe/customers/me/invoices", middlewares.LoadUserAccount, subscriptionEndpoint.GetCustomerInvoices)
	v1.Get("/stripe/customers/me/cards", middlewares.LoadUserAccount, subscriptionEndpoint.GetCustomerCards)
	v1.Delete("/stripe/cards", middlewares.LoadUserAccount, subscriptionEndpoint.RemoveCreditCard)
	v1.Put("/stripe/cards", middlewares.LoadUserAccount, subscriptionEndpoint.SetDefaultCreditCard)
	v1.Post("/stripe/create-setup-intent", middlewares.LoadUserAccount, subscriptionEndpoint.CreateSetupIntent)
	v1.Post("/stripe/create-customer-checkout-session", middlewares.LoadUserAccount, subscriptionEndpoint.CreateCustomerCheckoutSession)
	v1.Post("/stripe/create-customer-portal-session", middlewares.LoadUserAccount, subscriptionEndpoint.CreateCustomerPortalSession)
}
