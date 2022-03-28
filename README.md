# Starter SaaS Go API

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
invoice.payment_succeeded
invoice.payment_failed
customer.subscription.created
customer.subscription.updated
```

For local development, use the stripe-cli to build a local tunnel:

```bash
stripe listen --load-from-webhooks-api --forward-to localhost:3000
```

Finally, configure Stripe to retry failed payments for X days (https://dashboard.stripe.com/settings/billing/automatic Smart Retries section), and then cancel the subscription. 

Remember this value, it will be used in the `.env` file in `PAYMENT_FAILED_RETRY_DAYS` variable.

# Configuring .env

Below the meaning of every environment variable you can setup.


`PORT=":3000"`  the API server port number

`CORS_SITES="*"` allow only a specific domain to perform Ajax requests

`DATABASE="startersaas-db"` the MongoDB database 

`DATABASE_URI="mongodb://localhost:27017"` the MongoDB connection string

`JWT_SECRET="aaabbbccc"` set this value secret, very long and random

`JWT_EXPIRE="20"` # how many days the JWT token last

`FRONTEND_LOGIN_URL="http://localhost:5000/auth/login"` raplace http://localhost:5000 with the real production host of the React frontend

`MAILER_HOST='localhost'` the SMTP server host

`MAILER_PORT=1025` the SMTP server port

`MAILER_USERNAME='foo'` the SMTP server username

`MAILER_PASSWORD='bar'` the SMTP server password

`DEFAULT_EMAIL_FROM="noreply@startersaas.com"` send every notification email from this address

`LOCALE="en"` the default locale for registered users

`STRIPE_SECRET_KEY="sk_test_xyz"` the Stripe secret key

`FATTURA24_KEY="XYZ"` the Fattura 24 secret key (Italian market only)

`FATTURA24_URL="https://www.app.fattura24.com/api/v0.3/SaveDocument"` do not change this value

`TRIAL_DAYS=15` how many days a new user can work without subscribing

`PAYMENT_FAILED_RETRY_DAYS=7` how many days a user can work after the first failed payment (and before Stripe cancel the subscription)

`NOTIFIED_ADMIN_EMAIL="info@startersaas.com"` we notify this email when some events occur, like a new subscription, a failed payment and so on

`SIGNUP_WITH_ACTIVATE=true` set this value as true if you want to log the new registered user directly, without asking for email confirmation

`STARTER_PLAN_TYPE="starter"` set the plan to assign by default to a new customer 

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

* user registration of account with subdomain, email and password
* user email activation with 6 characters code and account creation
* resend activation code if not received
* user password reset through code sent by email
* user login
* user logout
* user change password once logged in
* account trial period
* edit of account billing informations
* subscription creation
* plan change
* add new credit card
* subscription cancel
* 3D Secure ready payments
* account's users list (by admins only)
* account's user create (by admins only)
* account's user update (by admins only)
* account's user delete (by admins only)

### API only

* stripe webhooks handling
* events notifications by email:
  - new user subscribed
  - successful payments
  - failed payments
* daily notifications by email:
  - expiring trials
  - failed payments
  - account suspension due to failed payments

### CREDITS

Author: Stefano Mancini <stefano.mancini@devinterface.com> 

Company: DevInterface SRL (https://www.devinterface.com)

Issues repository: https://github.com/devinterface/startersaas-issues

