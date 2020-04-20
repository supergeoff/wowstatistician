package auth

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/imroc/req"
)

var clientID = "7f5a41847fab493cb27e05fbcaecab0f"
var clientSecret = "wQLF9cATzaKjfJWwb5Jmuwrrz84VujcX"

// Token handle the response for CreateToken
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func encodeCred(id string, secret string) string {
	inputstr := fmt.Sprintf("%s:%s", id, secret)
	outputstr := base64.StdEncoding.EncodeToString([]byte(inputstr))
	return outputstr
}

// CreateToken generates a Bearer Token for API auth
func CreateToken() (string, error) {
	credentials := encodeCred(clientID, clientSecret)
	authstr := fmt.Sprintf("Basic %s", credentials)
	header := req.Header{
		"Authorization": authstr,
	}
	param := req.Param{
		"grant_type": "client_credentials",
	}
	request, err := req.Post("https://eu.battle.net/oauth/token", header, param)
	if err != nil {
		return "", errors.New("auth: could not generate token - " + err.Error())
	}
	var response Token
	request.ToJSON(&response)
	return response.AccessToken, nil
}
