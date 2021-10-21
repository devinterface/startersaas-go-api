package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"devinterface.com/startersaas-go-api/models"
	"github.com/Kamva/mgm/v3"
	strftime "github.com/jehiah/go-strftime"
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
	return success, err
}

// PaymentSucceeded function
func (webhookService *WebhookService) PaymentSucceeded(event stripe.Event) (success bool, err error) {
	sCustomerID := event.Data.Object["customer"]
	account, err := accountService.OneBy(bson.M{"stripeCustomerId": sCustomerID})
	if err != nil {
		return false, err
	}
	account.Active = true
	account.PaymentFailed = false
	account.PaymentFailedFirstAt = *new(time.Time)
	err = accountService.getCollection().Update(account)
	go emailService.SendStripeNotificationEmail(bson.M{"accountId": account.ID, "role": models.AdminRole}, "[Starter SAAS] Pagamento completato", "Il tuo abbonamento è stato rinnovato!")
	return err != nil, err
}

// PaymentFailed function
func (webhookService *WebhookService) PaymentFailed(event stripe.Event) (success bool, err error) {
	status := event.Data.Object["payment_intent"]
	if status != nil {
		return false, err
	}
	sCustomerID := event.Data.Object["customer"]
	account, _ := accountService.OneBy(bson.M{"stripeCustomerId": sCustomerID})
	account.PaymentFailed = true
	if account.PaymentFailedFirstAt.IsZero() {
		account.PaymentFailedFirstAt = time.Now()
	}
	err = accountService.getCollection().Update(account)
	paymentFailAfterDays, _ := strconv.Atoi(os.Getenv("PAYMENT_FAIL_AFTER_DAYS"))
	subscriptionDeactivatedAt := account.PaymentFailedFirstAt.AddDate(0, 0, paymentFailAfterDays)
	formattedSubscriptionDeactivatedAt := strftime.Format("%d/%m/%Y", subscriptionDeactivatedAt)

	go emailService.SendStripeNotificationEmail(bson.M{"accountId": account.ID, "role": models.AdminRole}, "[Starter SAAS] Pagamento fallito", fmt.Sprintf("Siamo spiacenti ma per qualche ragione il tuo pagamento non è andato a buon fine! Controlla le impostazioni della tua carta di credito. Il tuo account sarà sospeso il %s.", formattedSubscriptionDeactivatedAt))
	return err != nil, err
}

// TrialWillEnd function
func (webhookService *WebhookService) TrialWillEnd(event stripe.Event) (success bool, err error) {
	status := event.Data.Object["status"]
	if status != "active" {
		return false, err
	}
	sCustomerID := event.Data.Object["customer"]
	account, err := accountService.OneBy(bson.M{"stripeCustomerId": sCustomerID})
	go emailService.SendStripeNotificationEmail(bson.M{"accountId": account.ID, "role": models.AdminRole}, "[Starter SAAS] Il tuo abbonamento sta per scadere", "Il tuo abbonamento scadrà tra 3 giorni. Se non lo hai annullato, sarà rinnovato nuovamente.")
	return err != nil, err
}

// SubscriptionCreated function
func (webhookService *WebhookService) SubscriptionCreated(event stripe.Event) (success bool, err error) {
	status := event.Data.Object["status"]
	if status != "active" {
		return false, err
	}
	sCustomerID := event.Data.Object["customer"]
	account, err := accountService.OneBy(bson.M{"stripeCustomerId": sCustomerID})
	go emailService.SendStripeNotificationEmail(bson.M{"accountId": account.ID, "role": models.AdminRole}, "[Starter SAAS] Piano attivato", "Congratulazioni, il tuo piano è stato attivato!")
	return err != nil, err
}

// SubscriptionUpdated function
func (webhookService *WebhookService) SubscriptionUpdated(event stripe.Event) (success bool, err error) {
	status := event.Data.Object["status"]
	if status != "active" {
		return false, err
	}
	sCustomerID := event.Data.Object["customer"]
	account, err := accountService.OneBy(bson.M{"stripeCustomerId": sCustomerID})
	go emailService.SendStripeNotificationEmail(bson.M{"accountId": account.ID, "role": models.AdminRole}, "[Starter SAAS] Piano aggiornato", "Congratulazioni, il tuo piano è stato aggiornato!")
	return err != nil, err
}
