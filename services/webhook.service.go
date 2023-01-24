package services

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"devinterface.com/startersaas-go-api/models"
	"github.com/Kamva/mgm/v3"
	strftime "github.com/jehiah/go-strftime"
	"github.com/kataras/i18n"
	"github.com/stripe/stripe-go/v72"
	"github.com/thoas/go-funk"
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
	if event.Type == "invoice.paid" {
		success, err = webhookService.PaymentSuccessful(event)
	} else if event.Type == "invoice.payment_failed" {
		success, err = webhookService.PaymentFailed(event)
	} else if event.Type == "customer.subscription.created" {
		success, err = webhookService.SubscriptionUpdated(event)
	} else if event.Type == "customer.subscription.updated" {
		success, err = webhookService.SubscriptionUpdated(event)
	}
	return success, err
}

// PaymentSucceeded function
func (webhookService *WebhookService) PaymentSuccessful(event stripe.Event) (success bool, err error) {
	sCustomerID := event.Data.Object["customer"]
	account, err := accountService.OneBy(bson.M{"stripeCustomerId": sCustomerID})
	if err != nil {
		return false, err
	}
	account.PaymentFailed = false
	account.PaymentFailedFirstAt = *new(time.Time)
	account.PaymentFailedSubscriptionEndsAt = *new(time.Time)
	account.TrialPeriodEndsAt = *new(time.Time)
	_ = accountService.getCollection().Update(account)
	user, err := userService.OneBy(bson.M{"accountId": account.ID, "accountOwner": true})
	if err != nil {
		return false, err
	}

	billingReason := event.Data.Object["billing_reason"]
	if billingReason == "subscription_create" {
		go emailService.SendNotificationEmail(user.Email, i18n.Tr(user.Language, "webhookService.newSubscription.subject"), i18n.Tr(user.Language, "webhookService.newSubscription.message"), user.Language)
		go emailService.SendNotificationEmail(os.Getenv("NOTIFIED_ADMIN_EMAIL"), i18n.Tr(os.Getenv("LOCALE"), "webhookService.newSubscription.subject"), i18n.Tr(os.Getenv("LOCALE"), "webhookService.newSubscription.messageAdmin", map[string]interface{}{"Subdomain": account.Subdomain, "Email": user.Email}), os.Getenv("LOCALE"))

	}
	if billingReason == "subscription_update" {
		go emailService.SendNotificationEmail(user.Email, i18n.Tr(user.Language, "webhookService.subscriptionUpdated.subject"), i18n.Tr(user.Language, "webhookService.subscriptionUpdated.message"), user.Language)
		go emailService.SendNotificationEmail(os.Getenv("NOTIFIED_ADMIN_EMAIL"), i18n.Tr(os.Getenv("LOCALE"), "webhookService.subscriptionUpdated.subject"), i18n.Tr(os.Getenv("LOCALE"), "webhookService.subscriptionUpdated.messageAdmin", map[string]interface{}{"Subdomain": account.Subdomain, "Email": user.Email}), os.Getenv("LOCALE"))
	}
	if billingReason == "subscription_cycle" {
		go emailService.SendNotificationEmail(user.Email, i18n.Tr(user.Language, "webhookService.paymentSuccessful.subject"), i18n.Tr(user.Language, "webhookService.paymentSuccessful.message"), user.Language)
		go emailService.SendNotificationEmail(os.Getenv("NOTIFIED_ADMIN_EMAIL"), i18n.Tr(os.Getenv("LOCALE"), "webhookService.paymentSuccessful.subject"), i18n.Tr(os.Getenv("LOCALE"), "webhookService.paymentSuccessful.messageAdmin", map[string]interface{}{"Subdomain": account.Subdomain, "Email": user.Email}), os.Getenv("LOCALE"))
	}

	return err != nil, err
}

// PaymentFailed function
func (webhookService *WebhookService) PaymentFailed(event stripe.Event) (success bool, err error) {
	paymentIntent := event.Data.Object["payment_intent"]
	billingReason := event.Data.Object["billing_reason"]
	if paymentIntent != nil && billingReason != "subscription_update" {
		return false, err
	}
	sCustomerID := event.Data.Object["customer"]
	account, err := accountService.OneBy(bson.M{"stripeCustomerId": sCustomerID})
	if err != nil {
		return false, err
	}
	account.PaymentFailed = true
	if account.PaymentFailedFirstAt.IsZero() {
		account.PaymentFailedFirstAt = time.Now()
	}
	paymentFailedRetryDays, _ := strconv.Atoi(os.Getenv("PAYMENT_FAILED_RETRY_DAYS"))
	subscriptionDeactivatedAt := account.PaymentFailedFirstAt.AddDate(0, 0, paymentFailedRetryDays)
	account.PaymentFailedSubscriptionEndsAt = subscriptionDeactivatedAt
	_ = accountService.getCollection().Update(account)
	formattedSubscriptionDeactivatedAt := strftime.Format("%d/%m/%Y", subscriptionDeactivatedAt)
	user, err := userService.OneBy(bson.M{"accountId": account.ID, "accountOwner": true})
	if err != nil {
		return false, err
	}
	stripeHostedInvoiceUrl := event.Data.Object["hosted_invoice_url"]
	go emailService.SendNotificationEmail(user.Email, i18n.Tr(user.Language, "webhookService.paymentFailed.subject"), i18n.Tr(user.Language, "webhookService.paymentFailed.message", map[string]interface{}{"Date": formattedSubscriptionDeactivatedAt, "StripeHostedInvoiceUrl": stripeHostedInvoiceUrl}), user.Language)
	go emailService.SendNotificationEmail(os.Getenv("NOTIFIED_ADMIN_EMAIL"), i18n.Tr(os.Getenv("LOCALE"), "webhookService.paymentFailed.subject"), i18n.Tr(os.Getenv("LOCALE"), "webhookService.paymentFailed.messageAdmin", map[string]interface{}{"Subdomain": account.Subdomain, "Email": user.Email, "Date": formattedSubscriptionDeactivatedAt}), os.Getenv("LOCALE"))
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
	if err != nil {
		return false, err
	}
	sPlanMap := event.Data.Object["plan"].(map[string]interface{})
	if event.Data.Object["cancel_at"] != nil {
		sCancelAt := event.Data.Object["cancel_at"].(float64)
		account.SubscriptionExpiresAt = time.Unix(int64(sCancelAt), 0)
	} else {
		account.SubscriptionExpiresAt = *new(time.Time)
	}
	account.StripePlanID = sPlanMap["id"].(string)

	data, _ := ioutil.ReadFile("./stripe.conf.json")
	var payload interface{}
	json.Unmarshal(data, &payload)
	plans := payload.(map[string]interface{})
	found := funk.Filter(plans["plans"].([]interface{}), func(x interface{}) bool {
		i := x.(map[string]interface{})
		return i["id"] == account.StripePlanID
	}).([]interface{})
	plan := found[0].(map[string]interface{})
	account.PlanType = plan["planType"].(string)

	_ = accountService.getCollection().Update(account)
	return err != nil, err
}
