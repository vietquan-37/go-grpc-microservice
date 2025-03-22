package oauth2

const (
	GoogleTokenUrl    = "https://oauth2.googleapis.com/token"
	GoogleUserInfoUrl = "https://www.googleapis.com/oauth2/v1/userinfo"
)

type ExchangeTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}
type ExchangeTokenRequest struct {
	Code         string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	RedirectUri  string `json:"redirect_uri"`
}
type GoogleUserResponse struct {
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}
type ErrorResponse struct {
	Message string `json:"message"`
}
