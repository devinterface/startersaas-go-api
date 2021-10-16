package services

import (
	"errors"
	"os"

	"devinterface.com/startersaas-go-api/models"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/card"
	"github.com/stripe/stripe-go/v71/customer"
	"github.com/stripe/stripe-go/v71/invoice"
	"github.com/stripe/stripe-go/v71/paymentmethod"
	"github.com/stripe/stripe-go/v71/paymentsource"
	"github.com/stripe/stripe-go/v71/plan"
	"github.com/stripe/stripe-go/v71/sub"
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
func (subscriptionService *SubscriptionService) Subscribe(userID interface{}, planID string, sourceToken string) (subscription *stripe.Subscription, err error) {
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
	sCustomer, _ := customer.Get(account.StripeCustomerID, &stripe.CustomerParams{})

	noCrediCard := sCustomer.Sources.TotalCount == 0
	if noCrediCard {
		subscriptionService.AddCreditCard(account.ID, sourceToken)
	}

	sPlan, _ := plan.Get(planID, &stripe.PlanParams{})
	var activeSubscription *stripe.Subscription
	for _, s := range sCustomer.Subscriptions.Data {
		if s.Status == "active" {
			activeSubscription = s
		}
	}

	if activeSubscription != nil {
		params := &stripe.SubscriptionParams{CancelAtPeriodEnd: stripe.Bool(false), Plan: stripe.String(sPlan.ID)}
		subscription, err = sub.Update(activeSubscription.ID, params)
		if err != nil {
			return nil, err
		}
	} else {
		params := &stripe.SubscriptionParams{Customer: stripe.String(sCustomer.ID), Plan: stripe.String(sPlan.ID), PaymentBehavior: stripe.String("allow_incomplete")}
		params.AddExpand("latest_invoice.payment_intent")
		subscription, err = sub.New(params)
		if err != nil {
			return nil, err
		}
		account.FirstSubscription = false
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
	sCustomer, err = customer.Get(account.StripeCustomerID, &stripe.CustomerParams{})

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

	sCustomer, err = customer.Get(account.StripeCustomerID, &stripe.CustomerParams{})

	return sCustomer, err
}

// AddCreditCard function
func (subscriptionService *SubscriptionService) AddCreditCard(accountID interface{}, sourceToken string) (sCustomer *stripe.Customer, err error) {
	account := &models.Account{}
	err = accountService.getCollection().FindByID(accountID, account)
	if err != nil {
		return nil, err
	}
	if account.StripeCustomerID == "" {
		return nil, errors.New("user is not a stripe USER")
	}
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.CustomerSourceParams{
		Customer: stripe.String(account.StripeCustomerID),
		Source: &stripe.SourceParams{
			Token: stripe.String(sourceToken),
		},
	}

	_, err = paymentsource.New(params)

	if err != nil {
		return nil, err
	}

	customerParams := &stripe.CustomerParams{DefaultSource: stripe.String(sourceToken)}
	sCustomer, err = customer.Update(account.StripeCustomerID, customerParams)

	return sCustomer, err
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

	cardParams := &stripe.CardParams{
		Customer: stripe.String(account.StripeCustomerID),
	}
	_, err = card.Del(cardID, cardParams)

	if err != nil {
		return nil, err
	}
	sCustomer, err = customer.Get(account.StripeCustomerID, &stripe.CustomerParams{})

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

	customerParams := &stripe.CustomerParams{DefaultSource: stripe.String(cardID)}
	sCustomer, err = customer.Update(account.StripeCustomerID, customerParams)

	return sCustomer, err
}
