package main

import (
	"context"
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
	router.Route("GET", "users", getUsers)
	router.Route("GET", "users/:username", getUserCount)
	//router.Route("POST", "users/:username", createUser)
	//router.Route("POST", "users/:username/add", addToCount)
	//router.Route("POST", "users/:username/reset", resetCout)
}

func main() {
	lambda.Start(router.Handler)
}

func getUsers(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	response, err := counter.Users()
	if err != nil {
		return lmdrouter.HandleError(err)
	}
	return lmdrouter.MarshalResponse(http.StatusOK, nil, response)
}

func getUserCount(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	username := req.PathParameters["username"]
	response, err := counter.User(username)
	if err != nil {
		return lmdrouter.HandleError(err)
	}
	return lmdrouter.MarshalResponse(http.StatusOK, nil, response)

}
