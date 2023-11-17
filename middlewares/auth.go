package middlewares

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/alfredomagalhaes/fiber-user-api/controllers"
	"github.com/alfredomagalhaes/fiber-user-api/types"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
)

func AuthUser() fiber.Handler {
	return func(c *fiber.Ctx) error {

		issuerUrl := fmt.Sprintf(os.Getenv("OIDC_ISSUER_URL"), os.Getenv("AWS_REGION"), os.Getenv("COGNITO_USER_POOL_ID"))

		rawIdToken := c.GetReqHeaders()["Authorization"]
		rawIdToken = strings.ReplaceAll(rawIdToken, "Bearer ", "")
		requestPath := strings.ToLower(c.Path())
		requestMethod := string(c.Request().Header.Method())

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{
			Timeout:   time.Duration(6000) * time.Second,
			Transport: tr,
		}
		ctx := oidc.ClientContext(context.Background(), client)
		provider, err := oidc.NewProvider(ctx, issuerUrl)
		if err != nil {
			log.Printf("%v\n", err)
			return c.Status(http.StatusInternalServerError).JSON(controllers.ErrorResponse(errors.New("failed to connect with the authorization service")))
		}

		oidcConfig := &oidc.Config{
			ClientID: os.Getenv("OIDC_CLIENT_ID"),
		}
		verifier := provider.Verifier(oidcConfig)

		idToken, err := verifier.Verify(ctx, rawIdToken)
		if err != nil {
			log.Printf("%v\n", err)
			return c.Status(http.StatusUnauthorized).JSON(controllers.ErrorResponse(errors.New("invalid authentication token")))
		}

		var userClaims types.Claims

		err = idToken.Claims(&userClaims)

		if err != nil {
			log.Printf("%v\n", err)
			return c.Status(http.StatusInternalServerError).JSON(controllers.ErrorResponse(errors.New("failed to obtain access permissions")))
		}

		appRoles, err := getRoles() //TODO - remover esse tratamento da middleware

		if err != nil {
			log.Printf("%v\n", err)
			return c.Status(http.StatusInternalServerError).JSON(controllers.ErrorResponse(errors.New("failed to obtain access permissions")))
		}

		pathOk := false
		roleName := ""
		methodOk := false
		//Verifica se o usuário possui alguma role que permita
		//acessar a URL + Método executado
		for _, group := range userClaims.Groups {

			rolePath, ok := appRoles.RoleName[group]
			if ok {
				pathPattern := strings.ToLower(rolePath.Path)
				//Verifica através de uma expressão regex
				//se a url acessada pelo usuário ( 2º param ) bate
				//com o path da role ( 1º param )
				pathOk, _ = regexp.MatchString(pathPattern, requestPath)
				//Caso o path seja */* quer dizer que é uma role de admin
				if pathOk || pathPattern == "*/*" {
					roleName = group
					pathOk = true

				}
			}
			if pathOk {
				//Caso tenha encontrado uma role que permita acessar a role
				//verifica se é possível acessar o método na url
				methodOk = appRoles.RoleName[roleName].AllowedMethods[requestMethod]

			}

			if pathOk && methodOk {
				break
			}

		}

		if !pathOk {
			return c.Status(http.StatusUnauthorized).JSON(controllers.ErrorResponse(errors.New("usuário não tem permissão para acessar esse recurso")))
		}

		if !methodOk {
			return c.Status(http.StatusUnauthorized).JSON(controllers.ErrorResponse(errors.New("usuário não tem permissão para acessar esse recurso")))
		}

		return c.Next()
	}
}

func getRoles() (types.Role, error) {
	var roles types.Role

	rolesFile, err := os.Open("./seeds/roles.json")
	if err != nil {
		log.Printf("error while trying to read roles seed file\n%v", err)
		return roles, err
	}
	defer rolesFile.Close()

	parser := json.NewDecoder(rolesFile)
	err = parser.Decode(&roles)
	if err != nil {
		log.Printf("error while trying to decode file to json struct\n%v", err)
	}
	return roles, err
}
