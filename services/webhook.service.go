package services

import (
	"time"

	"devinterface.com/goaas-api-starter/models"
	"github.com/Kamva/mgm/v3"
	"github.com/stripe/stripe-go/v71"
	"go.mongodb.org/mongo-driver/bson"
)

// WebhookService struct
type WebhookService struct{ BaseService }

// HandleWebhook function
func (webhookService *WebhookService) HandleWebhook(payload map[string]interface{}, event stripe.Event) (success bool, err error) {
	go func(payload map[string]interface{}) {
		webhook := &models.Webhook{}
		webhook.Payload = payload
		mgm.CollectionByName("webhook").Create(webhook)
	}(payload)
	if event.Type == "invoice.payment_succeeded" {
		success, err = webhookService.PaymentSucceeded(event)
	} else if event.Type == "invoice.payment_failed" {
		success, err = webhookService.PaymentFailed(event)
	} else if event.Type == "customer.subscription.trial_will_end" {
		success, err = webhookService.TrialWillEnd(event)
	} else if event.Type == "customer.subscription.created" {
		success, err = webhookService.SubscriptionCreated(event)
	} else if event.Type == "customer.subscription.updated" {
		success, err = webhookService.SubscriptionUpdated(event)
	}
	return true, err
}

// PaymentSucceeded function
func (webhookService *WebhookService) PaymentSucceeded(event stripe.Event) (success bool, err error) {
	sCustomerID := event.Data.Object["customer"]
	account, err := accountService.OneBy(bson.M{"stripeCustomerId": sCustomerID})
	account.Active = true
	account.PaymentFailed = false
	periodEnd := event.Data.Object["period_end"].(float64)
	periodEndInt := int64(periodEnd)
	account.PeriodEndsAt = time.Unix(periodEndInt, 0)
	err = accountService.getCollection().Update(account)
	go emailService.SendStripeNotificationEmail(bson.M{"accountId": account.ID, "role": models.AdminRole}, "[MiniMarket24] Pagamento completato", "Il tuo abbonamento è stato rinnovato!")
	return err != nil, err
}

// PaymentFailed function
func (webhookService *WebhookService) PaymentFailed(event stripe.Event) (success bool, err error) {
	sCustomerID := event.Data.Object["customer"]
	account, err := accountService.OneBy(bson.M{"stripeCustomerId": sCustomerID})
	account.Active = true
	account.PaymentFailed = true
	err = accountService.getCollection().Update(account)
	go emailService.SendStripeNotificationEmail(bson.M{"accountId": account.ID, "role": models.AdminRole}, "[MiniMarket24] Pagamento fallito", "Il tuo abbonamento è stato rinnovato! Controlla le impostazioni della tua carta di credito.")
	return err != nil, err
}

// TrialWillEnd function
func (webhookService *WebhookService) TrialWillEnd(event stripe.Event) (success bool, err error) {
	sCustomerID := event.Data.Object["customer"]
	account, err := accountService.OneBy(bson.M{"stripeCustomerId": sCustomerID})
	go emailService.SendStripeNotificationEmail(bson.M{"accountId": account.ID, "role": models.AdminRole}, "[MiniMarket24] Il tuo abbonamento sta per scadere", "Il tuo abbonamento scadrà tra 3 giorni. Se non lo hai annullato, sarà rinnovato nuovamente.")
	return err != nil, err
}

// SubscriptionCreated function
func (webhookService *WebhookService) SubscriptionCreated(event stripe.Event) (success bool, err error) {
	sCustomerID := event.Data.Object["customer"]
	account, err := accountService.OneBy(bson.M{"stripeCustomerId": sCustomerID})
	go emailService.SendStripeNotificationEmail(bson.M{"accountId": account.ID, "role": models.AdminRole}, "[MiniMarket24] Piano attivato", "Congratulazioni, il tuo piano è stato attivato!")
	return err != nil, err
}

// SubscriptionUpdated function
func (webhookService *WebhookService) SubscriptionUpdated(event stripe.Event) (success bool, err error) {
	sCustomerID := event.Data.Object["customer"]
	account, err := accountService.OneBy(bson.M{"stripeCustomerId": sCustomerID})
	go emailService.SendStripeNotificationEmail(bson.M{"accountId": account.ID, "role": models.AdminRole}, "[MiniMarket24] Piano aggiornato", "Congratulazioni, il tuo piano è stato aggiornato!")
	return err != nil, err
}
