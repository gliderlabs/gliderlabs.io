package auth0

import (
	"net/http"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/pkg/auth0"
	"golang.org/x/oauth2"
)

type UserInfo auth0.UserInfo
type User auth0.User

type LogoutFn func(http.ResponseWriter, *http.Request) error
type LoginFn func(http.ResponseWriter, *http.Request, *oauth2.Token) error

var LogoutCallback LogoutFn
var LoginCallback LoginFn

func Client() *auth0.Client {
	return &auth0.Client{
		ClientID:     com.GetString("client_id"),
		ClientSecret: com.GetString("client_secret"),
		Domain:       com.GetString("domain"),
		CallbackURL:  com.GetString("callback_url"),
		Token:        com.GetString("api_token"),
	}
}
