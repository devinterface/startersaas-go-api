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
	account.PaymentFailedSubscriptionEndsAt = *new(time.Time)
	account.TrialPeriodEndsAt = *new(time.Time)
	err = accountService.getCollection().Update(account)
	user, _ := userService.OneBy(bson.M{"accountId": account.ID})
	go emailService.SendNotificationEmail(user.Email, "[Starter SAAS] Payment completed", "Congratulations, your subscription has been renewed.")
	go emailService.SendNotificationEmail(os.Getenv("NOTIFIED_ADMIN_EMAIL"), "[Starter SAAS] Payment completed", fmt.Sprintf("%s - %s - paid a subscription", account.Subdomain, user.Email))
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
	paymentFailedRetryDays, _ := strconv.Atoi(os.Getenv("PAYMENT_FAILED_RETRY_DAYS"))
	subscriptionDeactivatedAt := account.PaymentFailedFirstAt.AddDate(0, 0, paymentFailedRetryDays)
	account.PaymentFailedSubscriptionEndsAt = subscriptionDeactivatedAt

	err = accountService.getCollection().Update(account)
	formattedSubscriptionDeactivatedAt := strftime.Format("%d/%m/%Y", subscriptionDeactivatedAt)

	user, _ := userService.OneBy(bson.M{"accountId": account.ID})
	go emailService.SendNotificationEmail(user.Email, "[Starter SAAS] Payment failed", fmt.Sprintf("Your payment wasn't successful. Please check your payment card and retry. Your subscription will be deactivated on %s.", formattedSubscriptionDeactivatedAt))
	go emailService.SendNotificationEmail(os.Getenv("NOTIFIED_ADMIN_EMAIL"), "[Starter SAAS] Payment failed", fmt.Sprintf("%s - %s - has a failed payment. His subscription will be deactivated on %s.", account.Subdomain, user.Email, formattedSubscriptionDeactivatedAt))
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
	user, _ := userService.OneBy(bson.M{"accountId": account.ID})
	go emailService.SendNotificationEmail(user.Email, "[Starter SAAS] New subscription activated", "Congratulations, your subscription has been activated.")
	go emailService.SendNotificationEmail(os.Getenv("NOTIFIED_ADMIN_EMAIL"), "[Starter SAAS] New subscription activated", fmt.Sprintf("%s - %s - activated a subscription.", account.Subdomain, user.Email))
	return err != nil, err
}

// SubscriptionUpdated function
func (webhookService *WebhookService) SubscriptionUpdated(event stripe.Event) (success bool, err error) {
	status := event.Data.Object["status"]
	if status != "active" {
		return false, err
	}
	sCustomerID := event.Data.Object["customer"]
	account, _ := accountService.OneBy(bson.M{"stripeCustomerId": sCustomerID})
	sPlanMap := event.Data.Object["plan"].(map[string]interface{})

	if event.Data.Object["cancel_at"] != nil {
		sCancelAt := event.Data.Object["cancel_at"].(float64)
		account.SubscriptionExpiresAt = time.Unix(int64(sCancelAt), 0)
	} else {
		account.SubscriptionExpiresAt = *new(time.Time)
	}
	account.StripePlanID = sPlanMap["id"].(string)
	err = accountService.getCollection().Update(account)
	user, _ := userService.OneBy(bson.M{"accountId": account.ID})
	go emailService.SendNotificationEmail(user.Email, "[Starter SAAS] Subscription updated", "Congratulations, your subscription has been updated.")
	go emailService.SendNotificationEmail(os.Getenv("NOTIFIED_ADMIN_EMAIL"), "[Starter SAAS] Subscription updated", fmt.Sprintf("%s - %s - updated a subscription.", account.Subdomain, user.Email))
	return err != nil, err
}
