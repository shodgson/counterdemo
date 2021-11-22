package counter

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
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

func TestAddUser(t *testing.T) {
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

	u, err = Increment(id, 1)
	assert.Equal(t, 1, u.Count)
	_, err = Increment(id, 1)
	_, err = Increment(id, 1)
	u, err = User(id)
	assert.Equal(t, 3, u.Count)

}

/*
func TestValidIncrements(t *testing.T) {
	tests := []struct {
		name string,
		startCount int,
		premium bool,
		add int,
		valid bool
	}{
		{}
	}
	u, := CreateUser("valid-increment-1", "valid-increment)
}
*/

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func loadEnvVariables() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func init() {
	loadEnvVariables()
	SetupTable()
	err := Delete(testUser.Id)
	if err != nil {
		log.Fatal("Unable to delete user")
	}
}
