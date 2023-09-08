# StarterSaaS Go API

This project contains everything you need to setup a fully featured SaaS API in 5 minutes.

# Installation

Copy `.env.example` into `.env` and `stripe.conf.json.example` into `stripe.conf.json`.

Create a startersaas newtwork typing:

```bash
docker network create startersaas-network
```

Then build the containers

```bash
docker compose build
```

And finally, run the application

```bash
docker compose up
```

Application will be reachable on

```bash
http://localhost:3000
```

# Stripe setup

Configure your stripe webhooks by setting as "Endpoint URL" :

```
https://<my_startersaas_api_domain>/api/v1/stripe/webhook
```

and events below:

```
invoice.paid
invoice.payment_failed
customer.subscription.created
customer.subscription.updated
```

For local development, use the stripe-cli to build a local tunnel:

```bash
stripe listen --forward-to localhost:3000/api/v1/stripe/webhook
```

Finally, configure Stripe to retry failed payments for X days (https://dashboard.stripe.com/settings/billing/automatic Smart Retries section), and then cancel the subscription.

Remember this value, it will be used in the `.env` file in `PAYMENT_FAILED_RETRY_DAYS` variable.

# Update dependencies

Rin this command:

```
docker compose run --rm startersaas-go-api go get -u && go mod tidy
```

# Configuring .env

Below the meaning of every environment variable you can setup.

`PORT=":3000"` the API server port number

`CORS_SITES="*"` allow only a specific domain to perform Ajax requests

`DATABASE="startersaas-db"` the MongoDB database

`DATABASE_URI="mongodb://localhost:27017"` the MongoDB connection string

`JWT_SECRET="aaabbbccc"` set this value secret, very long and random

`JWT_EXPIRE=20` # how many days the JWT token last

`FRONTEND_LOGIN_URL="http://localhost:5000/auth/login"` raplace http://localhost:5000 with the real production host of the React frontend

`MAILER_HOST='localhost'` the SMTP server host

`MAILER_PORT=1025` the SMTP server port

`MAILER_USERNAME='foo'` the SMTP server username

`MAILER_PASSWORD='bar'` the SMTP server password

`MAILER_SSL=false`, true if the SMTP server uses SSL

`DEFAULT_EMAIL_FROM="noreply@startersaas.com"` send every notification email from this address

`LOCALE="en"` the default locale for registered users

`STRIPE_SECRET_KEY="sk_test_xyz"` the Stripe secret key

`TRIAL_DAYS=15` how many days a new user can work without subscribing

`PAYMENT_FAILED_RETRY_DAYS=7` how many days a user can work after the first failed payment (and before Stripe cancel the subscription)

`NOTIFIED_ADMIN_EMAIL="info@startersaas.com"` we notify this email when some events occur, like a new subscription, a failed payment and so on

`SIGNUP_WITH_ACTIVATE=true` set this value as true if you want to log the new registered user directly, without asking for email confirmation

`STARTER_PLAN_TYPE="starter"` set the plan to assign by default to a new customer. Must match one of the plans defined in `stripe.conf.json`

`FRONTEND_CUSTOMER_PORTAL_REDIRECT_URL="http://localhost:3010/dashboard"` the URL to forward after actions on Stripe Customer Portal

### Docker variables

`APP_PORT=3000` the port the API is available in local machine

`MONGO_PORT=27017` the port the Database is available in local machine

`REDIS_PORT=6379` the port the Redis server is available in local machine

`MAILHOG_UI_PORT=8025` the port the Mailhog is available in local machine

# Configuring stripe.conf.json

In this file you have to add the stripe API public key, in the `publicKey` field

Then for every product you want to sell, copy it's price_id (usually starts with price_xxx) and paste it in the "id" key.

```
{
  "id": "price_XYZ",
  "title": "Starter",
  "price": 4.90,
  "currency": "EUR",
  "features": [
    "1 project",
    "0 tags",
    "star entries",
    "1 user",
    "3 days data retention",
    "no push notifications"
  ],
  "monthly": true,
  "planType": "starter"
}
```

Then sets its title, its price (in cents, the same you have configured in Stripe) and the list of features you want to show in the frontend pricing table.

Set `"monthly":true` if your plan is billed on monthly basis, otherwise we consider it billed yearly.

Set `"planType"` with your plan code to a more user friendly knowledge of the current plan.

# Features

### API and Frontend

- [x] user registration of account with subdomain, email and password
- [x] user email activation with 6 characters code and account creation
- [x] resend activation code if not received
- [x] user password reset through code sent by email
- [x] user login
- [x] user logout
- [x] user change password once logged in
- [x] account trial period
- [x] edit of account billing information
- [x] subscription creation
- [x] plan change
- [x] add new credit card
- [x] remove credit card
- [x] subscription cancel
- [x] subscription re enable
- [x] 3D Secure ready payments
- [x] subscription handling via Stripe customer portal
- [x] account's users list (by admins only)
- [x] account's user create (by admins only)
- [x] account's user update (by admins only)
- [x] account's user delete (by admins only)
- [x] teams handling

### API only

- [x] stripe webhooks handling
- events notifications by email:
  - [x] new user subscribed
  - [x] successful payments
  - [x] failed payments
- daily notifications by email:
  - [x] expiring trials
  - [x] failed payments
  - [x] account suspension due to failed payments

### TODO

- [ ] signup with Google

### CREDITS

Author: Stefano Mancini <stefano.mancini@devinterface.com>

Company: DevInterface SRL (https://www.devinterface.com)

### License

Licensed under the [MIT License](https://github.com/devinterface/startersaas-go-api/blob/master/LICENSE)
