package stripe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaymentURL(t *testing.T) {
	url, err := PaymentURL("test@email.com", "id123")
	assert.Nil(t, err)
	checkoutBase := "https://checkout.stripe.com/"
	assert.Equal(t, checkoutBase, url[0:len(checkoutBase)])
}
