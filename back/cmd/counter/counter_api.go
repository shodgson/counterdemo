package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/guregu/dynamo"
	lmdrouter "github.com/shodgson/lmdrouterv2"
	"github.com/shodgson/speedbuild/internal/apilogger"
	"github.com/shodgson/speedbuild/internal/counter"
	"github.com/shodgson/speedbuild/internal/stripe"
)

var router *lmdrouter.Router
var db *dynamo.DB

func init() {
	router = lmdrouter.NewRouter("", apilogger.Logger)
	router.Route(http.MethodGet, "users", getUsers)
	router.Route(http.MethodGet, "users/:username", getUserCount)
	router.Route(http.MethodPatch, "count", setCount)
	router.Route(http.MethodGet, "subscription/activation_url", activationUrl)
	router.Route(http.MethodPost, "subscription/webhook", stripeWebhook)
	router.Route(http.MethodPost, "subscription/cancel", cancelStripe)
}

func main() {
	lambda.Start(router.Handler)
}

func response(obj interface{}, err error) (events.APIGatewayV2HTTPResponse, error) {
	if err != nil {
		if err == dynamo.ErrNotFound {
			return lmdrouter.MarshalResponse(http.StatusNotFound, nil, "Item not found")
		}
		if err == counter.ErrAccessDenied {
			return lmdrouter.MarshalResponse(http.StatusForbidden, nil, "User does not have permission to perform this action")
		}
		return lmdrouter.HandleError(err)
	}
	return lmdrouter.MarshalResponse(http.StatusOK, nil, obj)
}

func getUsers(ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (
	events.APIGatewayV2HTTPResponse,
	error,
) {
	return response(counter.Users())
}

func getUserCount(ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (
	events.APIGatewayV2HTTPResponse,
	error,
) {
	log.Println(req.PathParameters)
	id := req.PathParameters["username"]
	return response(counter.User(id))
}

func setCount(ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (
	events.APIGatewayV2HTTPResponse,
	error,
) {
	id := req.RequestContext.Authorizer.JWT.Claims["sub"]
	body := struct {
		Add   int  `json:"add"`
		Reset bool `json:"reset",omitempty`
	}{}
	err := json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		response(nil, err)
	}
	if body.Reset {
		return response(counter.Reset(id))
	}
	return response(counter.Increment(id, body.Add))
}

func activationUrl(ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (
	events.APIGatewayV2HTTPResponse,
	error,
) {
	email := req.RequestContext.Authorizer.JWT.Claims["email"]
	id := req.RequestContext.Authorizer.JWT.Claims["sub"]
	url, err := stripe.PaymentURL(email, id)
	respBody := struct {
		URL string `json:"URL"`
	}{
		URL: url,
	}
	return response(respBody, err)
}

func stripeWebhook(ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (
	events.APIGatewayV2HTTPResponse,
	error,
) {
	stripeSignature := req.Headers["stripe-signature"]
	result, userId, subId, err := stripe.ProcessWebhook([]byte(req.Body), stripeSignature)
	if err != nil {
		return response(nil, err)
	}
	log.Println(result)
	switch result {
	case stripe.HookPaid:
		_, err = counter.SetAccount(userId, true, subId)
	case stripe.HookCancelled:
		log.Println("HookCancelled")
		log.Println(userId)
		_, err = counter.SetAccount(userId, false, "")
	default:
		err = errors.New("Unknown webhook")
	}
	return response("", err)
}

func cancelStripe(ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (
	events.APIGatewayV2HTTPResponse,
	error,
) {
	id := req.RequestContext.Authorizer.JWT.Claims["sub"]
	u, err := counter.User(id)
	if err != nil {
		return response("", err)
	}
	err = stripe.Unsubscribe(u.StripeSubscriptionID, true, false)
	return response(nil, err)
}
