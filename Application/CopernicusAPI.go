package Application

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

const openIDTokenPath = "/auth/realms/CDSE/protocol/openid-connect/token"

var client = resty.New()

func init() {

	client.SetBaseURL("https://identity.dataspace.copernicus.eu")

}

func GetRefreshToken(credentials Credentials) (RefreshToken, error) {
	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type": "password",
			"client_id":  "cdse-public",
			"username":   credentials.Username,
			"password":   credentials.Password,
		}).
		Post(openIDTokenPath)

	if err != nil {
		return RefreshToken{}, err
	}

	if !resp.IsSuccess() {
		return RefreshToken{}, fmt.Errorf("error: %s", resp.Status())
	}

	var refreshToken RefreshToken
	err = json.Unmarshal(resp.Body(), &refreshToken)
	if err != nil {
		return RefreshToken{}, err
	}

	return refreshToken, nil
}

func GetAccessToken(refreshToken RefreshToken) (AccessToken, error) {
	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "refresh_token",
			"client_id":     "cdse-public",
			"refresh_token": refreshToken.Token,
		}).
		Post(openIDTokenPath)

	if !resp.IsSuccess() {
		return AccessToken{}, fmt.Errorf("error: %s", resp.Status())
	}

	if err != nil {
		return AccessToken{}, err
	}

	var accessToken AccessToken
	err = json.Unmarshal(resp.Body(), &accessToken)
	if err != nil {
		return AccessToken{}, err
	}

	return accessToken, nil
}
