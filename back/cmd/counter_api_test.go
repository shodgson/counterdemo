// +build integration

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/shodgson/speedbuild/internal/cognito"
	"github.com/shodgson/speedbuild/internal/counter"
	"github.com/stretchr/testify/assert"
)

var apiUrl string

var testUsers = []counter.CountItem{
	{
		Name:  "apitestuser12@example.com",
		Count: 0,
	},
	{
		Name:  "intruder1@example.com",
		Count: 0,
	},
}

var testCognitoUsers = []cognito.User{
	{
		Username: testUsers[0].Name,
		Password: "testpassword123",
	},
	{
		Username: testUsers[1].Name,
		Password: "l33tpassword",
	},
}

func TestAPI(t *testing.T) {

	setupTest()

	// Create cognito users
	err := testCognitoUsers[0].SignUp()
	assert.Nil(t, err)
	err = testCognitoUsers[1].SignUp()
	assert.Nil(t, err)
	type CognitoClaims struct {
		Username string `json:"username"`
		jwt.StandardClaims
	}

	// User creation
	for i, u := range testCognitoUsers {
		token, _ := jwt.Parse(*u.AccessToken, nil)
		testUsers[i].Id = fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["username"])
		testUser := testUsers[i]
		checkEndpoint(t, "GET", "users", nil, http.StatusOK)
		checkEndpoint(t, "GET", "users/"+testUser.Id, nil, http.StatusNotFound)
		checkEndpoint(t, "POST", "users/"+testUser.Name, nil, http.StatusUnauthorized)
		checkEndpoint(t, "POST", "users/"+testUser.Name, u.AccessToken, http.StatusOK)
		checkEndpoint(t, "GET", "users/"+testUser.Id, nil, http.StatusOK)
	}

	// Add
	addTest := []struct {
		user       int
		authorized bool
		add        int
		resultCode int
		resultSum  int
	}{
		{0, false, 1, http.StatusUnauthorized, 0},
		{0, true, 1, http.StatusOK, 1},
		{1, true, 1, http.StatusOK, 1},
		{0, true, 1, http.StatusOK, 2},
		{0, true, 3, http.StatusOK, 5},
		{0, true, 1, http.StatusForbidden, 2},
	}

	for _, aT := range addTest {
		requestBody := []byte(fmt.Sprintf("{\"add\": %d}", aT.add))
		user := testUsers[aT.user]
		token := testCognitoUsers[aT.user].AccessToken
		if !aT.authorized {
			token = nil
		}
		resp, err := sendRequest("PATCH", "count", token, requestBody)
		assert.Nil(t, err)
		assert.Equal(t, aT.resultCode, resp.StatusCode)
		if aT.resultCode == http.StatusOK {
			defer resp.Body.Close()
			rbody, _ := ioutil.ReadAll(resp.Body)
			assert.Equal(t, fmt.Sprintf(`{"id":"%s","name":"%s","count":%d}`, user.Id, user.Name, aT.resultSum), string(rbody))
		}
	}

	teardown()

}

func checkEndpoint(t *testing.T, method string, path string, token *string, expectedStatus int) {
	resp, err := sendRequest(method, path, token, nil)
	assert.Nil(t, err, "%s %s - request failed", method, path)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, expectedStatus, resp.StatusCode, "Response body: %s\n%s %s - response: %d (expected %d)", string(body), method, path, resp.StatusCode, expectedStatus)
}

func sendRequest(method string, path string, token *string, body []byte) (*http.Response, error) {
	requestBody := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, apiUrl+"/"+path, requestBody)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	if token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
	}
	client := &http.Client{}
	return client.Do(req)
}

func setupTest() {
	// Set environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiUrl = os.Getenv("API_URL")
	counter.SetupTable()
	cognito.SetupCognito()

	// Clean up test data
	err = counter.Delete(testUsers[0].Name)
	if err != nil {
		log.Fatal("Error deleting test user 0", err)
	}
	err = counter.Delete(testUsers[1].Name)
	if err != nil {
		log.Fatal("Error deleting test user 1", err)
	}
}

func teardown() {
	err := testCognitoUsers[0].Delete()
	if err != nil {
		log.Fatal("Error deleting cognito test user", err)
	}
	err = testCognitoUsers[1].Delete()
	if err != nil {
		log.Fatal("Error deleting cognito test user", err)
	}
}
