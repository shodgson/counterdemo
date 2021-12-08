package stripe

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	stripe "github.com/stripe/stripe-go/v72"
	stripesession "github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/stripe/stripe-go/v72/webhook"
)

var priceId string
var stripeKey string
var baseURL string
var endpointSecret string

const stripeSuccessPath = "?payment=success"
const stripeCancelPath = "?payment=cancel"

func SetupConfiguration() {
	priceId = os.Getenv("STRIPE_PRICE_ID")
	stripeKey = os.Getenv("STRIPE_KEY")
	baseURL = os.Getenv("BASE_URL")
	endpointSecret = os.Getenv("STRIPE_ENDPOINT_SECRET")
}

func init() {
	SetupConfiguration()
}

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

	log.Println("HOOK!")
	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "checkout.session.completed":
		var checkout stripe.CheckoutSession
		err = json.Unmarshal(event.Data.Raw, &checkout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			return
		}
		userId = checkout.ClientReferenceID
		result = HookPaid
		subscriptionId = checkout.Subscription.ID

	case "customer.subscription.updated":
		var subscription stripe.Subscription
		err = json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			return
		}
		subscriptionId = subscription.ID
		// TODO: if payment fails, find customer and mark as inactive

	case "customer.subscription.deleted":
		log.Println("deleted")
		var subscription stripe.Subscription
		err = json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			return
		}
		result = HookCancelled
		subscriptionId = subscription.ID
		userId = subscription.Metadata["id"]
		return

	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}
	return
}
