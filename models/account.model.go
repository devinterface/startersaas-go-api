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
	CompanyCountry                  string    `json:"companyCountry" bson:"companyCountry"`
	StripeCustomerID                string    `json:"stripeCustomerId" bson:"stripeCustomerId"`
	PaymentFailed                   bool      `json:"paymentFailed" bson:"paymentFailed"`
	PaymentFailedFirstAt            time.Time `json:"paymentFailedFirstAt" bson:"paymentFailedFirstAt"`
	PaymentFailedSubscriptionEndsAt time.Time `json:"paymentFailedSubscriptionEndsAt" bson:"paymentFailedSubscriptionEndsAt"`
	PrivacyAccepted                 bool      `json:"privacyAccepted" bson:"privacyAccepted"`
	MarketingAccepted               bool      `json:"marketingAccepted" bson:"marketingAccepted"`
	TrialPeriodEndsAt               time.Time `json:"trialPeriodEndsAt" bson:"trialPeriodEndsAt"`
	StripePlanID                    string    `json:"stripePlanId" bson:"stripePlanId"`
	SubscriptionExpiresAt           time.Time `json:"subcriptionExpiresAt" bson:"subcriptionExpiresAt"`
	PlanType                        string    `json:"planType" bson:"planType"`
}

// AccountSerializer function
type AccountSerializer struct {
	*structomap.Base
}

// ShowAccountSerializer function
func ShowAccountSerializer() *AccountSerializer {
	a := &AccountSerializer{structomap.New()}
	a.UseCamelCase().Pick("ID", "Subdomain", "CompanyName", "CompanyVat", "CompanyBillingAddress", "CompanySdi", "CompanyPec", "CompanyPhone", "CompanyEmail", "CompanyCountry", "PaymentFailed", "PaymentFailedFirstAt", "TrialPeriodEndsAt", "PaymentFailedSubscriptionEndsAt", "PrivacyAccepted", "MarketingAccepted", "StripePlanID", "SubscriptionExpiresAt", "PlanType", "CreatedAt", "UpdatedAt").AddFunc("SubscriptionStatus", func(a interface{}) interface{} {
		account := a.(*Account)
		return account.SubscriptionStatus()
	})
	return a
}

const (
	SubscriptionTrial         = "trial"
	SubscriptionPaymentFailed = "payment_failed"
	SubscriptionDeactivated   = "deactivated"
	SubscriptionActive        = "active"
)

func (account *Account) SubscriptionStatus() string {
	if account.TrialPeriodEndsAt.After(time.Now()) {
		return SubscriptionTrial
	} else if account.TrialPeriodEndsAt != *new(time.Time) && account.TrialPeriodEndsAt.Before(time.Now()) {
		return SubscriptionDeactivated
	} else if account.PaymentFailed && account.PaymentFailedSubscriptionEndsAt != *new(time.Time) && account.PaymentFailedSubscriptionEndsAt.After(time.Now()) {
		return SubscriptionPaymentFailed
	} else if account.PaymentFailed && account.PaymentFailedSubscriptionEndsAt != *new(time.Time) && account.PaymentFailedSubscriptionEndsAt.Before(time.Now()) {
		return SubscriptionDeactivated
	} else if account.SubscriptionExpiresAt != *new(time.Time) && account.SubscriptionExpiresAt.Before(time.Now()) {
		return SubscriptionDeactivated
	} else {
		return SubscriptionActive
	}
}

const (
	StarterPlan = "starter"
	BasicPlan   = "basic"
	PremiumPlan = "premium"
)
