package counter

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testUser = CountItem{
	Id:    "abc-123",
	Name:  "test",
	Count: 0,
}

func TestUsers(t *testing.T) {
	_, err := Users()
	assert.Nil(t, err)
}

func TestCounter(t *testing.T) {
	id := testUser.Id
	name := testUser.Name
	users, _ := Users()
	numUsers := len(users)

	// Create user
	u, err := CreateUser(id, name)
	assert.Nil(t, err)
	assert.Equal(t, name, u.Name)

	users, _ = Users()
	assert.Equal(t, numUsers+1, len(users))
	assert.Equal(t, 0, u.Count)

	// Get user
	u, err = User(id)
	assert.Equal(t, name, u.Name)
	assert.Equal(t, 0, u.Count)

	// Increment
	u, err = Increment(id, 1)
	assert.Nil(t, err)
	assert.Equal(t, 1, u.Count)
	_, err = Increment(id, 1)
	assert.Nil(t, err)
	_, err = Increment(id, 1)
	assert.Nil(t, err)
	u, err = User(id)
	assert.Equal(t, 3, u.Count)

	// Premium limits (free account)
	_, err = Increment(id, 3)
	assert.NotNil(t, err)
	u, err = Increment(id, 2)
	assert.Equal(t, 5, u.Count)
	u, err = Increment(id, 1)
	assert.NotNil(t, err)

	// Account management
	subID := "sub123"
	u, err = SetAccount(id, true, subID)
	assert.Nil(t, err)
	assert.Equal(t, true, u.Premium)
	assert.Equal(t, subID, u.StripeSubscriptionID)
	u, err = SetAccount(id, false, "")
	assert.Nil(t, err)
	assert.Equal(t, false, u.Premium)
	u, err = SetAccount(id, true, subID)
	assert.Nil(t, err)
	assert.Equal(t, true, u.Premium)

	// Premium limits (premium account)
	u, err = Increment(id, 1)
	assert.Equal(t, 6, u.Count)

	// Reset
	u, err = Reset(id)
	assert.Equal(t, 0, u.Count)

}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func init() {
	err := Delete(testUser.Id)
	if err != nil {
		log.Fatal("Unable to delete user")
	}
}
