package stripe

import (
	"encoding/json"
	"fmt"
	"os"

	stripe "github.com/stripe/stripe-go/v72"
	stripesession "github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/stripe/stripe-go/v72/webhook"
)

var priceId = os.Getenv("STRIPE_PRICE_ID")
var stripeKey = os.Getenv("STRIPE_KEY")
var baseURL = os.Getenv("BASE_URL")
var endpointSecret = os.Getenv("STRIPE_ENDPOINT_SECRET")

const stripeSuccessPath = "?payment=success"
const stripeCancelPath = "?payment=cancel"

func Unsubscribe(subscriptionID string, now bool, undo bool) (err error) {
	stripe.Key = stripeKey
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(!undo),
	}

	if now {
		_, err = sub.Cancel(subscriptionID, nil)
	} else {
		_, err = sub.Update(subscriptionID, params)
	}

	return
}

func PaymentURL(email string, id string) (url string, err error) {
	stripe.Key = stripeKey
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(baseURL + stripeSuccessPath),
		CancelURL:  stripe.String(baseURL + stripeCancelPath),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		ClientReferenceID: stripe.String(id),
		CustomerEmail:     stripe.String(email),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"id": id,
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Price:    stripe.String(priceId),
				Quantity: stripe.Int64(1),
			},
		},
	}
	s, err := stripesession.New(params)
	if err != nil {
		return "", err
	}
	return s.URL, err
}

const (
	HookPaid = iota
	HookCancelled
)

func ProcessWebhook(body []byte, signature string) (result int, userId string, subscriptionId string, err error) {
	stripe.Key = stripeKey

	event, err := webhook.ConstructEvent(body, signature, endpointSecret)

	if err != nil {
		fmt.Printf("sig: %s, secret: %s, key: %s\n", signature, endpointSecret, stripe.Key)
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		return
	}

	var subscription stripe.Subscription
	err = json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
		return
	}
	subscriptionId = subscription.ID
	userId = subscription.Metadata["id"]

	switch event.Type {
	case "customer.subscription.created":
		result = HookPaid
		return
	case "customer.subscription.deleted":
		result = HookCancelled
		return
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
		result = -1
	}
	return
}
