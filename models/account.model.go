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
	FirstSubscription               bool      `json:"firstSubscription" bson:"firstSubscription"`
	PaymentFailed                   bool      `json:"paymentFailed" bson:"paymentFailed"`
	PaymentFailedFirstAt            time.Time `json:"paymentFailedFirstAt" bson:"paymentFailedFirstAt"`
	PaymentFailedSubscriptionEndsAt time.Time `json:"paymentFailedSubscriptionEndsAt" bson:"paymentFailedSubscriptionEndsAt"`
	PrivacyAccepted                 bool      `json:"privacyAccepted" bson:"privacyAccepted"`
	MarketingAccepted               bool      `json:"marketingAccepted" bson:"marketingAccepted"`
	TrialPeriodEndsAt               time.Time `json:"trialPeriodEndsAt" bson:"trialPeriodEndsAt"`
}

// AccountSerializer function
type AccountSerializer struct {
	*structomap.Base
}

// ShowAccountSerializer function
func ShowAccountSerializer() *AccountSerializer {
	a := &AccountSerializer{structomap.New()}
	a.UseCamelCase().Pick("ID", "Subdomain", "CompanyName", "CompanyVat", "CompanyBillingAddress", "CompanySdi", "CompanyPec", "CompanyPhone", "CompanyEmail", "Active", "FirstSubscription", "PaymentFailed", "PaymentFailedFirstAt", "TrialPeriodEndsAt", "PaymentFailedSubscriptionEndsAt", "PrivacyAccepted", "MarketingAccepted", "CreatedAt", "UpdatedAt")
	return a
}
