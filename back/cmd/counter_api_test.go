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

var testUser = counter.CountItem{
	Name:  "apitestuser12@example.com",
	Count: 0,
}

var testCognitoUser = cognito.User{
	Username: testUser.Name,
	Password: "testpassword123",
}

func TestAPI(t *testing.T) {

	setupTest()

	// Create cognito user
	err := testCognitoUser.SignUp()
	assert.Nil(t, err)
	type CognitoClaims struct {
		Username string `json:"username"`
		jwt.StandardClaims
	}
	token, _ := jwt.Parse(*testCognitoUser.AccessToken, nil)
	testUser.Id = fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["username"])

	// User
	checkEndpoint(t, "GET", "users", nil, http.StatusOK)
	checkEndpoint(t, "GET", "users/"+testUser.Id, nil, http.StatusNotFound)
	checkEndpoint(t, "POST", "users/"+testUser.Name, nil, http.StatusUnauthorized)
	checkEndpoint(t, "POST", "users/"+testUser.Name, testCognitoUser.AccessToken, http.StatusOK)
	checkEndpoint(t, "GET", "users/"+testUser.Id, nil, http.StatusOK)

	// Add 1
	resp, err := sendRequest("PATCH", "users/"+testUser.Id, nil, []byte(`{"add": 1}`))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	resp, err = sendRequest("PATCH", "users/"+testUser.Name, testCognitoUser.AccessToken, []byte(`{"add": 1}`))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
	rbody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, fmt.Sprintf(`{"id":"%s","name":"%s","count":1}`, testUser.Id, testUser.Name), string(rbody))

	// Add 2 more
	resp, err = sendRequest("PATCH", "users/"+testUser.Name, testCognitoUser.AccessToken, []byte(`{"add": 2}`))
	defer resp.Body.Close()
	rbody, _ = ioutil.ReadAll(resp.Body)
	assert.Equal(t, fmt.Sprintf(`{"id":"%s","name":"%s","count":3}`, testUser.Id, testUser.Name), string(rbody))

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
	err = counter.Delete(testUser.Name)
	if err != nil {
		log.Fatal("Error deleting test user", err)
	}
}

func teardown() {
	err := testCognitoUser.Delete()
	if err != nil {
		log.Fatal("Error deleting cognito test user", err)
	}
}
