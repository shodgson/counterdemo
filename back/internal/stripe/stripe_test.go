package stripe

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestPaymentURL(t *testing.T) {
	url, err := PaymentURL("test@email.com")
	assert.Nil(t, err)
	checkoutBase := "https://checkout.stripe.com/"
	assert.Equal(t, checkoutBase, url[0:len(checkoutBase)])
}

func TestWebhook(t *testing.T) {

}

func init() {
	loadEnvVariables()
	SetupConfiguration()
}

func loadEnvVariables() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
