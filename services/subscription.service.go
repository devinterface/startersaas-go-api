package services

import (
	"errors"
	"math"
	"os"
	"time"

	"devinterface.com/startersaas-go-api/models"
	"github.com/Kamva/mgm/v3/operator"
	strftime "github.com/jehiah/go-strftime"
	"github.com/kataras/i18n"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/billingportal/session"
	checkout_session "github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/plan"
	"github.com/stripe/stripe-go/v72/setupintent"
	"github.com/stripe/stripe-go/v72/sub"
	"go.mongodb.org/mongo-driver/bson"
)

// SubscriptionService struct
type SubscriptionService struct{ BaseService }

// CreateCustomer function
func (subscriptionService *SubscriptionService) CreateCustomer(userID interface{}) (sCustomer *stripe.Customer, err error) {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	user := &models.User{}
	err = userService.getCollection().FindByID(userID, user)
	if err != nil {
		return nil, err
	}
	account := &models.Account{}
	err = accountService.getCollection().FindByID(user.AccountID, account)
	if err != nil {
		return nil, err
	}

	params := &stripe.CustomerParams{
		Name:  stripe.String(account.CompanyName),
		Email: stripe.String(user.Email),
		Address: &stripe.AddressParams{
			Line1:   stripe.String(account.CompanyBillingAddress),
			City:    stripe.String(account.CompanyCountry),
			Country: stripe.String(account.CompanyCountry),
		},
	}
	params.AddMetadata("companyName", account.CompanyName)
	params.AddMetadata("address", account.CompanyBillingAddress)
	params.AddMetadata("vat", account.CompanyVat)
	params.AddMetadata("subdomain", account.Subdomain)
	params.AddMetadata("sdi", account.CompanySdi)
	params.AddMetadata("phone", account.CompanyPhone)
	params.AddMetadata("userName", user.Name)
	params.AddMetadata("userSurname", user.Surname)

	sCustomer, _ = customer.New(params)
	account.StripeCustomerID = sCustomer.ID
	err = accountService.getCollection().Update(account)
	return sCustomer, err
}

// Subscribe function
func (subscriptionService *SubscriptionService) Subscribe(userID interface{}, planID string) (subscription *stripe.Subscription, err error) {
	user := &models.User{}
	err = userService.getCollection().FindByID(userID, user)
	if err != nil {
		return nil, err
	}
	account := &models.Account{}
	err = accountService.getCollection().FindByID(user.AccountID, account)
	if err != nil {
		return nil, err
	}
	if account.StripeCustomerID == "" {
		_, err = subscriptionService.CreateCustomer(userID)
		if err != nil {
			return nil, err
		}
		err = accountService.getCollection().FindByID(user.AccountID, account)
		if err != nil {
			return nil, err
		}
	}
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	customerParams := &stripe.CustomerParams{}
	customerParams.AddExpand("subscriptions")
	sCustomer, _ := customer.Get(account.StripeCustomerID, customerParams)

	sPlan, _ := plan.Get(planID, &stripe.PlanParams{})
	var activeSubscription *stripe.Subscription
	for _, s := range sCustomer.Subscriptions.Data {
		if s.Status == "active" {
			activeSubscription = s
		}
	}

	if activeSubscription != nil {
		params := &stripe.SubscriptionParams{CancelAtPeriodEnd: stripe.Bool(false),
			Plan:              stripe.String(sPlan.ID),
			ProrationBehavior: stripe.String("always_invoice"),
			PaymentBehavior:   stripe.String("default_incomplete"),
			AutomaticTax: &stripe.SubscriptionAutomaticTaxParams{
				Enabled: stripe.Bool(true),
			}}
		params.AddExpand("latest_invoice.payment_intent")
		subscription, err = sub.Update(activeSubscription.ID, params)
		if err != nil {
			return nil, err
		}
	} else {

		for _, s := range sCustomer.Subscriptions.Data {
			if s.Status != "active" {
				sub.Cancel(s.ID, nil)
			}
		}

		params := &stripe.SubscriptionParams{Customer: stripe.String(sCustomer.ID),
			Plan:            stripe.String(sPlan.ID),
			PaymentBehavior: stripe.String("default_incomplete"),
			AutomaticTax: &stripe.SubscriptionAutomaticTaxParams{
				Enabled: stripe.Bool(true),
			}}
		params.AddExpand("latest_invoice.payment_intent")
		subscription, err = sub.New(params)
		if err != nil {
			return nil, err
		}
		err = accountService.getCollection().Update(account)
	}

	return subscription, err
}

// GetCustomer function
func (subscriptionService *SubscriptionService) GetCustomer(accountID interface{}) (sCustomer *stripe.Customer, err error) {
	account := &models.Account{}
	err = accountService.getCollection().FindByID(accountID, account)
	if err != nil {
		return nil, err
	}
	if account.StripeCustomerID == "" {
		return nil, errors.New("user is not a stripe USER")
	}
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	customerParams := &stripe.CustomerParams{}
	customerParams.AddExpand("subscriptions")
	sCustomer, err = customer.Get(account.StripeCustomerID, customerParams)

	return sCustomer, err
}

// GetCustomerInvoices function
func (subscriptionService *SubscriptionService) GetCustomerInvoices(accountID interface{}) (invoices []stripe.Invoice, err error) {
	account := &models.Account{}
	err = accountService.getCollection().FindByID(accountID, account)
	if err != nil {
		return nil, err
	}
	if account.StripeCustomerID == "" {
		return nil, errors.New("user is not a stripe USER")
	}
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	params := &stripe.InvoiceListParams{Customer: stripe.String(account.StripeCustomerID)}
	i := invoice.List(params)
	for i.Next() {
		in := i.Invoice()
		invoices = append(invoices, *in)
	}
	return invoices, err
}

// GetCustomerCards function
func (subscriptionService *SubscriptionService) GetCustomerCards(accountID interface{}) (cards []stripe.PaymentMethod, err error) {
	account := &models.Account{}
	err = accountService.getCollection().FindByID(accountID, account)
	if err != nil {
		return nil, err
	}
	if account.StripeCustomerID == "" {
		return nil, errors.New("user is not a stripe USER")
	}
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(account.StripeCustomerID),
		Type:     stripe.String("card"),
	}
	i := paymentmethod.List(params)
	for i.Next() {
		card := i.PaymentMethod()
		cards = append(cards, *card)
	}
	return cards, err
}

// CancelSubscription function
func (subscriptionService *SubscriptionService) CancelSubscription(accountID interface{}, subscriptionId string) (sCustomer *stripe.Customer, err error) {
	account := &models.Account{}
	err = accountService.getCollection().FindByID(accountID, account)
	if err != nil {
		return nil, err
	}
	if account.StripeCustomerID == "" {
		return nil, errors.New("user is not a stripe USER")
	}
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.SubscriptionParams{CancelAtPeriodEnd: stripe.Bool(true)}
	sub.Update(subscriptionId, params)

	customerParams := &stripe.CustomerParams{}
	customerParams.AddExpand("subscriptions")
	sCustomer, err = customer.Get(account.StripeCustomerID, customerParams)

	return sCustomer, err
}

// CreateSetupIntent function
func (subscriptionService *SubscriptionService) CreateSetupIntent(accountID interface{}) (setupIntent *stripe.SetupIntent, err error) {
	account := &models.Account{}
	err = accountService.getCollection().FindByID(accountID, account)
	if err != nil {
		return nil, err
	}
	if account.StripeCustomerID == "" {
		return nil, errors.New("user is not a stripe USER")
	}
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.SetupIntentParams{
		Customer: &account.StripeCustomerID,
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
	}
	setupIntent, _ = setupintent.New(params)

	return setupIntent, err
}

// RemoveCreditCard function
func (subscriptionService *SubscriptionService) RemoveCreditCard(accountID interface{}, cardID string) (sCustomer *stripe.Customer, err error) {
	account := &models.Account{}
	err = accountService.getCollection().FindByID(accountID, account)
	if err != nil {
		return nil, err
	}
	if account.StripeCustomerID == "" {
		return nil, errors.New("user is not a stripe USER")
	}
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	paymentmethod.Detach(
		cardID,
		nil,
	)

	if err != nil {
		return nil, err
	}
	customerParams := &stripe.CustomerParams{}
	customerParams.AddExpand("subscriptions")
	sCustomer, err = customer.Get(account.StripeCustomerID, customerParams)

	return sCustomer, err
}

// SetDefaultCreditCard function
func (subscriptionService *SubscriptionService) SetDefaultCreditCard(accountID interface{}, cardID string) (sCustomer *stripe.Customer, err error) {
	account := &models.Account{}
	err = accountService.getCollection().FindByID(accountID, account)
	if err != nil {
		return nil, err
	}
	if account.StripeCustomerID == "" {
		return nil, errors.New("user is not a stripe USER")
	}
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	customerInvoiceSettingsParams := &stripe.CustomerInvoiceSettingsParams{DefaultPaymentMethod: stripe.String(cardID)}
	customerParams := &stripe.CustomerParams{InvoiceSettings: customerInvoiceSettingsParams}
	sCustomer, err = customer.Update(account.StripeCustomerID, customerParams)

	return sCustomer, err
}

// RunNotifyExpiringTrials function
func (subscriptionService *SubscriptionService) RunNotifyExpiringTrials() (err error) {
	params := bson.M{"trialPeriodEndsAt": bson.M{operator.Lt: time.Now().AddDate(0, 0, 3), operator.Gt: time.Now()}}
	accounts, err := accountService.FindBy(params)
	for _, account := range accounts {
		user, _ := userService.OneBy(bson.M{"accountId": account.ID, "accountOwner": true})
		daysToExpire := int(math.Round(time.Until(account.TrialPeriodEndsAt).Hours() / 24))
		go emailService.SendNotificationEmail(user.Email, i18n.Tr(user.Language, "subscriptionService.runNotifyExpiringTrials.subject", map[string]interface{}{"DaysToExpire": daysToExpire}), i18n.Tr(user.Language, "subscriptionService.runNotifyExpiringTrials.message", map[string]interface{}{"DaysToExpire": daysToExpire}), user.Language)
	}
	return err
}

// RunNotifyPaymentFailed function
func (subscriptionService *SubscriptionService) RunNotifyPaymentFailed() (err error) {
	params := bson.M{"paymentFailed": true, "paymentFailedSubscriptionEndsAt": bson.M{operator.Lt: time.Now().AddDate(0, 0, 3), operator.Gt: time.Now()}}
	accounts, err := accountService.FindBy(params)
	for _, account := range accounts {
		user, _ := userService.OneBy(bson.M{"accountId": account.ID, "accountOwner": true})
		formattedPaymentFailedSubscriptionEndsAt := strftime.Format("%d/%m/%Y", account.PaymentFailedSubscriptionEndsAt)
		daysToExpire := int(math.Round(time.Until(account.PaymentFailedSubscriptionEndsAt).Hours() / 24))
		go emailService.SendNotificationEmail(user.Email, i18n.Tr(user.Language, "subscriptionService.runNotifyPaymentFailed.subject", map[string]interface{}{"DaysToExpire": daysToExpire}), i18n.Tr(user.Language, "subscriptionService.runNotifyPaymentFailed.message", map[string]interface{}{"Date": formattedPaymentFailedSubscriptionEndsAt}), user.Language)
	}
	return err
}

// CreateCustomerCheckoutSession function
func (subscriptionService *SubscriptionService) CreateCustomerCheckoutSession(userID interface{}, planID string) (redirect_url string, err error) {
	user := &models.User{}
	err = userService.getCollection().FindByID(userID, user)
	if err != nil {
		return "", err
	}
	account := &models.Account{}
	err = accountService.getCollection().FindByID(user.AccountID, account)
	if err != nil {
		return "", err
	}
	if account.StripeCustomerID == "" {
		_, err = subscriptionService.CreateCustomer(userID)
		if err != nil {
			return "", err
		}
		err = accountService.getCollection().FindByID(user.AccountID, account)
		if err != nil {
			return "", err
		}
	}
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	params := &stripe.CheckoutSessionParams{
		SuccessURL:   stripe.String(os.Getenv("FRONTEND_CUSTOMER_PORTAL_REDIRECT_URL")),
		CancelURL:    stripe.String(os.Getenv("FRONTEND_CUSTOMER_PORTAL_REDIRECT_URL")),
		Mode:         stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		Customer:     stripe.String(account.StripeCustomerID),
		AutomaticTax: &stripe.CheckoutSessionAutomaticTaxParams{Enabled: stripe.Bool(true)},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Price: stripe.String(planID),
				// For metered billing, do not pass quantity
				Quantity: stripe.Int64(1),
			},
		},
	}

	s, _ := checkout_session.New(params)
	return s.URL, err
}

// CreateCustomerPortalSession function
func (subscriptionService *SubscriptionService) CreateCustomerPortalSession(AccountID interface{}) (redirect_url string, err error) {
	account := &models.Account{}
	err = accountService.getCollection().FindByID(AccountID, account)
	if err != nil {
		return "", err
	}
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(account.StripeCustomerID),
		ReturnURL: stripe.String(os.Getenv("FRONTEND_CUSTOMER_PORTAL_REDIRECT_URL")),
	}
	s, _ := session.New(params)
	return s.URL, err
}
