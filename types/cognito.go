package types

import cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

type CognitoConfig struct {
	Client       *cognito.CognitoIdentityProvider
	UserPoolID   string
	ClientID     string
	ClientSecret string
	Token        string
}
