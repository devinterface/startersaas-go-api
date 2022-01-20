package models

import (
	"time"

	"github.com/Kamva/mgm/v3"
	"github.com/devinterface/structomap"
)

// Account struct
type Account struct {
	mgm.DefaultModel                `bson:",inline"`
	Subdomain                       string    `json:"subdomain" bson:"subdomain"`
	CompanyName                     string    `json:"companyName" bson:"companyName"`
	CompanyVat                      string    `json:"companyVat" bson:"companyVat"`
	CompanyBillingAddress           string    `json:"companyBillingAddress" bson:"companyBillingAddress"`
	CompanySdi                      string    `json:"companySdi" bson:"companySdi"`
	CompanyPec                      string    `json:"companyPec" bson:"companyPec"`
	CompanyPhone                    string    `json:"companyPhone" bson:"companyPhone"`
	CompanyEmail                    string    `json:"companyEmail" bson:"companyEmail"`
	StripeCustomerID                string    `json:"stripeCustomerId" bson:"stripeCustomerId"`
	Active                          bool      `json:"active" bson:"active"`
	PaymentFailed                   bool      `json:"paymentFailed" bson:"paymentFailed"`
	PaymentFailedFirstAt            time.Time `json:"paymentFailedFirstAt" bson:"paymentFailedFirstAt"`
	PaymentFailedSubscriptionEndsAt time.Time `json:"paymentFailedSubscriptionEndsAt" bson:"paymentFailedSubscriptionEndsAt"`
	PrivacyAccepted                 bool      `json:"privacyAccepted" bson:"privacyAccepted"`
	MarketingAccepted               bool      `json:"marketingAccepted" bson:"marketingAccepted"`
	TrialPeriodEndsAt               time.Time `json:"trialPeriodEndsAt" bson:"trialPeriodEndsAt"`
	StripePlanID                    string    `json:"stripePlanId" bson:"stripePlanId"`
	SubscriptionExpiresAt           time.Time `json:"subcriptionExpiresAt" bson:"subcriptionExpiresAt"`
}

// AccountSerializer function
type AccountSerializer struct {
	*structomap.Base
}

// ShowAccountSerializer function
func ShowAccountSerializer() *AccountSerializer {
	a := &AccountSerializer{structomap.New()}
	a.UseCamelCase().Pick("ID", "Subdomain", "CompanyName", "CompanyVat", "CompanyBillingAddress", "CompanySdi", "CompanyPec", "CompanyPhone", "CompanyEmail", "Active", "PaymentFailed", "PaymentFailedFirstAt", "TrialPeriodEndsAt", "PaymentFailedSubscriptionEndsAt", "PrivacyAccepted", "MarketingAccepted", "StripePlanID", "SubscriptionExpiresAt", "CreatedAt", "UpdatedAt")
	return a
}

const (
	SubscriptionPendingActivation = "pending_activation"
	SubscriptionTrial             = "trial"
	SubscriptionPaymentFailed     = "payment_failed"
	SubscriptionDeactivated       = "deactivated"
	SubscriptionActive            = "active"
)

func (account *Account) SubscriptionStatus() string {
	if !account.Active {
		return SubscriptionPendingActivation
	}
	if account.TrialPeriodEndsAt.After(time.Now()) {
		return SubscriptionTrial
	}
	if account.PaymentFailed && account.PaymentFailedSubscriptionEndsAt.After(time.Now()) {
		return SubscriptionPaymentFailed
	}
	if account.PaymentFailed && account.PaymentFailedSubscriptionEndsAt.Before(time.Now()) {
		return SubscriptionDeactivated
	}
	if account.SubscriptionExpiresAt.Before(time.Now()) {
		return SubscriptionDeactivated
	}
	return SubscriptionActive
}
