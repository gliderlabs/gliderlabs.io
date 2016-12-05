package auth0

import "github.com/gliderlabs/comlab/pkg/com"

func init() {
	com.Register("auth0", &Component{},
		com.Option("client_id", "", "Auth0 client ID"),
		com.Option("client_secret", "", "Auth0 client secret"),
		com.Option("domain", "", "Auth0 domain"),
		com.Option("callback_url", "/_auth/callback", "Auth0 callback URL"),
		com.Option("logout_url", "/_auth/logout", "URL to wrap Auth0 logout"),
		com.Option("api_token", "", "Auth0 API bearer token"))
}

// Component ...
type Component struct{}
