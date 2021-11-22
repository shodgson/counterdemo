package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/guregu/dynamo"
	lmdrouter "github.com/shodgson/lmdrouterv2"
	"github.com/shodgson/speedbuild/internal/apilogger"
	"github.com/shodgson/speedbuild/internal/counter"
)

var router *lmdrouter.Router
var db *dynamo.DB

func init() {
	router = lmdrouter.NewRouter("", apilogger.Logger)
	router.Route(http.MethodGet, "users", getUsers)
	router.Route(http.MethodPost, "users/:username", createUser)
	router.Route(http.MethodGet, "users/:username", getUserCount)
	router.Route(http.MethodPatch, "users/:username", setCount)
	//router.Route("POST", "users/:username", playing())
	//router.Route("POST", "users/:username/add", addToCount)
	//router.Route("POST", "users/:username/reset", resetCout)
}

func main() {
	lambda.Start(router.Handler)
}

func response(obj interface{}, err error) (events.APIGatewayV2HTTPResponse, error) {
	if err != nil {
		if err == dynamo.ErrNotFound {
			return lmdrouter.MarshalResponse(http.StatusNotFound, nil, "Item not found")
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

func createUser(ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (
	events.APIGatewayV2HTTPResponse,
	error,
) {
	username := req.PathParameters["username"]
	id := req.RequestContext.Authorizer.JWT.Claims["username"]
	return response(counter.CreateUser(id, username))
}

func setCount(ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (
	events.APIGatewayV2HTTPResponse,
	error,
) {
	id := req.RequestContext.Authorizer.JWT.Claims["username"]
	body := struct {
		Add   int  `json:"add"`
		Reset bool `json:"reset"`
	}{}
	err := json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		response(nil, err)
	}
	return response(counter.Increment(id, body.Add))
}
