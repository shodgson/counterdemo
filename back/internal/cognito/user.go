package cognito

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var cip *cognito.CognitoIdentityProvider
var cognitoClientId string

func SetupCognito() {
	region := os.Getenv("AWS_REGION")
	cognitoClientId = os.Getenv("COGNITO_CLIENT_ID")
	cip = cognito.New(session.New(), aws.NewConfig().WithRegion(region))
}

type User struct {
	Username          string
	Password          string
	Email             string
	IdToken           *string
	AccessToken       *string
	PreferredUsername string
}

func (u *User) SignUp() error {
	cognitoUser := &cognito.SignUpInput{
		Username: aws.String(u.Username),
		Password: aws.String(u.Password),
		ClientId: aws.String(cognitoClientId),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(u.Email),
			},
		},
	}
	_, err := cip.SignUp(cognitoUser)
	if err != nil {
		return err
	}
	return u.GetTokens()
}

func (u *User) GetTokens() error {

	authTry := &cognito.InitiateAuthInput{
		//AuthFlow: aws.String("USER_SRP_AUTH"),
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(u.Username),
			"PASSWORD": aws.String(u.Password),
		},
		ClientId: aws.String(cognitoClientId),
	}

	initRes, err := cip.InitiateAuth(authTry)
	if err != nil {
		fmt.Println("Could not authenticate:", err)
		return err
	}
	u.IdToken = initRes.AuthenticationResult.IdToken
	u.AccessToken = initRes.AuthenticationResult.AccessToken
	return err

}

func (u *User) SetUsername() error {
	attributes := []*cognito.AttributeType{
		{
			Name:  aws.String("preferred_username"),
			Value: aws.String(u.PreferredUsername),
		},
	}
	input := &cognito.UpdateUserAttributesInput{
		AccessToken: u.AccessToken,
		//ClientMetadata: map[string]*string{
		//"": nil,
		//},
		UserAttributes: attributes,
	}
	_, err := cip.UpdateUserAttributes(input)
	if err != nil {
		return err
	}
	return u.GetTokens()

}

func (u *User) Delete() error {
	cognitoUser := &cognito.DeleteUserInput{
		AccessToken: u.AccessToken,
	}
	_, err := cip.DeleteUser(cognitoUser)
	return err
}

func init() {
	SetupCognito()
}
