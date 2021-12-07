package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shodgson/speedbuild/internal/counter"
)

func handler(event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	fmt.Printf("PreSignup of user: %s\n", event.UserName)
	event.Response.AutoConfirmUser = true
	event.Response.AutoVerifyEmail = true
	_, err := counter.CreateUser(event.UserName, event.Request.UserAttributes["email"])
	if err != nil {
		fmt.Printf("Error creating database entry: %v\n", err)
	}
	return event, nil
}

func main() {
	lambda.Start(handler)
}
