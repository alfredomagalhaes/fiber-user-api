package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/alfredomagalhaes/fiber-user-api/types"
	"github.com/aws/aws-sdk-go/aws"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gofiber/fiber/v2"
)

var errProcReq error = errors.New("error while processing the request")
var ErrBodyParse error = errors.New("failed to parse the body from request")

func LoginUser(cngCfg types.CognitoConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(types.LoginRequest)

		if err := c.BodyParser(&request); err != nil {
			log.Printf("error while processing the request: %v", err)
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse(ErrBodyParse))
		}

		params := map[string]*string{
			"USERNAME": aws.String(request.UserName),
			"PASSWORD": aws.String(request.Password),
		}

		secretHash := computeSecretHash(os.Getenv("OIDC_CLIENT_SECRET"), request.UserName, os.Getenv("OIDC_CLIENT_ID"))
		params["SECRET_HASH"] = aws.String(secretHash)

		authTry := &cognito.InitiateAuthInput{
			AuthFlow: aws.String("USER_PASSWORD_AUTH"),
			AuthParameters: map[string]*string{
				"USERNAME":    aws.String(*params["USERNAME"]),
				"PASSWORD":    aws.String(*params["PASSWORD"]),
				"SECRET_HASH": aws.String(*params["SECRET_HASH"]),
			},
			ClientId: aws.String(os.Getenv("OIDC_CLIENT_ID")), // this is the app client ID
		}
		authResp, err := cngCfg.Client.InitiateAuth(authTry)
		if err != nil {
			log.Printf("error while authenticating with cognito\n%v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"data": authResp})
		}

		cngCfg.Token = *authResp.AuthenticationResult.AccessToken
		return c.Status(http.StatusOK).JSON(fiber.Map{"data": authResp})
	}
}

func computeSecretHash(clientSecret string, username string, clientId string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

/*
func LoginUser(c *fiber.Ctx) error {

	request := new(types.LoginRequest)

	if err := c.BodyParser(&request); err != nil {
		log.Printf("error while processing the request: %v", err)
		return c.Status(http.StatusBadRequest).JSON(ErrorResponse(ErrBodyParse))
	}

	urlData := url.Values{}
	urlData.Set("grant_type", "password")
	urlData.Set("client_id", os.Getenv("OIDC_CLIENT_ID"))
	urlData.Set("client_secret", os.Getenv("OIDC_CLIENT_SECRET"))
	urlData.Set("username", request.UserName)
	urlData.Set("password", request.Password)
	urlData.Set("scope", "profile email")
	//urlData.Add("scope", "email")
	encodedData := urlData.Encode()

	req, err := http.NewRequest("POST", os.Getenv("OIDC_TOKEN_URL"), strings.NewReader(encodedData))
	if err != nil {
		log.Printf("Got error %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse(errProcReq))
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(urlData.Encode())))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Got error %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse(errProcReq))
	}
	responseData, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		log.Printf("Got error %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse(errProcReq))
	}

	var oAuth2Token oauth2.Token

	err = json.Unmarshal(responseData, &oAuth2Token)
	if err != nil {
		log.Printf("Got error %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse(errProcReq))
	}

	return c.Status(response.StatusCode).JSON(oAuth2Token)
	//return c.Status(http.StatusOK).JSON(fiber.Map{"success": true, "message": "user logged in", "token": "123456"})
}
*/
// Creates a fiber Map object to
// standardize errors responses
func ErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"success": false,
		"data":    "",
		"error":   err.Error(),
	}
}
