package oauth2

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func ExchangeToken(request ExchangeTokenRequest) (rsp *ExchangeTokenResponse, err error) {
	data := url.Values{}
	data.Set("client_id", request.ClientId)
	data.Set("client_secret", request.ClientSecret)
	data.Set("grant_type", request.GrantType)
	data.Set("redirect_uri", request.RedirectUri)
	data.Set("code", request.Code)
	req, err := http.NewRequest(http.MethodPost, GoogleTokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {

		return nil, errors.New(string(body))
	}
	if err = json.Unmarshal(body, &rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}
func GetUserInfo(alt, token string) (rsp *GoogleUserResponse, err error) {
	params := url.Values{}
	params.Set("access_token", token)
	params.Set("alt", alt)
	urlWithParams := GoogleUserInfoUrl + "?" + params.Encode()
	req, err := http.NewRequest(http.MethodGet, urlWithParams, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &rsp); err != nil {
		return nil, err
	}
	return rsp, nil

}
