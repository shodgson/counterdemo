package counter

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

var tableName = "counts"
var sess = session.Must(session.NewSession())
var db = dynamo.New(sess, &aws.Config{
	Region: aws.String(os.Getenv("AWS_REGION")),
})
var stage = os.Getenv("SLS_STAGE")
var countTable = db.Table(fmt.Sprintf("%s-%s", stage, tableName))

// Database entry
type CountItem struct {
	Username string `dynamo:"username" json:"username"`
	Count    int    `dynamo:"count" json:"count"`
}

func Users() (results []CountItem, err error) {
	err = countTable.Scan().All(&results)
	return
}

func User(username string) (c CountItem, err error) {
	err = countTable.Get("username", username).One(&c)
	return
}
