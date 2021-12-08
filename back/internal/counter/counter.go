package counter

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

// Configuration
var tableName = "counts"
var FreeLimit = 5

var sess = session.Must(session.NewSession(&aws.Config{
	Region: aws.String(os.Getenv("AWS_REGION")),
}))
var db = dynamo.New(sess, nil)
var stage = os.Getenv("STAGE")
var countTable = db.Table(fmt.Sprintf("%s-%s", stage, tableName))

var ErrAccessDenied = errors.New("Access denied")

// Database entry
type CountItem struct {
	Id    string `dynamo:"id" json:"id"`
	Name  string `dynamo:"name" json:"name"`
	Count int    `dynamo:"count" json:"count"`

	// Premium account
	StripeSubscriptionID string `dynamo:"stripe_subscription_id" json:"-"`
	Premium              bool   `dynamo:"premium" json:"premium,omitempty"`
}

// Methods
func Users() (results []CountItem, err error) {
	err = countTable.Scan().All(&results)
	return
}

func User(id string) (c CountItem, err error) {
	err = countTable.Get("id", id).One(&c)
	return
}

func CreateUser(id string, name string) (c CountItem, err error) {
	c = CountItem{Name: name, Id: id}
	err = countTable.Put(c).If("attribute_not_exists(id)").Run()
	return
}

func Increment(id string, value int) (c CountItem, err error) {
	err = countTable.Update("id", id).
		Add("count", value).
		If("$ = ? OR $ <= ?", "premium", true, "count", FreeLimit-value).
		Value(&c)
	if ae, ok := err.(awserr.RequestFailure); ok {
		if ae.Code() == "ConditionalCheckFailedException" {
			err = ErrAccessDenied
		}
	}
	return
}

func Reset(id string) (c CountItem, err error) {
	err = countTable.Update("id", id).
		Set("count", 0).
		Value(&c)
	if ae, ok := err.(awserr.RequestFailure); ok {
		if ae.Code() == "ConditionalCheckFailedException" {
			err = ErrAccessDenied
		}
	}
	return
}

func SetAccount(id string, premium bool, subscriptionID string) (c CountItem, err error) {
	err = countTable.Update("id", id).
		Set("premium", premium).
		Set("stripe_subscription_id", subscriptionID).
		Value(&c)
	return
}

func Delete(id string) error {
	return countTable.Delete("id", id).Run()
}
