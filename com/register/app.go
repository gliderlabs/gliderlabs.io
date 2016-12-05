package register

import (
	"encoding/gob"
	"net/http"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/gliderlabs.io/com/auth0"
	"github.com/gliderlabs/gliderlabs.io/com/web"
	auth0lib "github.com/gliderlabs/pkg/auth0"
	"github.com/stripe/stripe-go"
	"golang.org/x/oauth2"
)

func init() {
	gob.Register(auth0lib.UserInfo{})
	gob.Register(map[string]interface{}{})
}

func (c *Component) AppPreStart() error {
	auth0.LoginCallback = func(w http.ResponseWriter, r *http.Request, token *oauth2.Token) error {
		userinfo, err := auth0.Client().UserInfo(token)
		if err != nil {
			return err
		}

		web.SessionSet(r, w, "_access_token", token.AccessToken)
		web.SessionSet(r, w, "_auth_id", userinfo["user_id"].(string))
		web.SessionSet(r, w, "user_name", userinfo["name"].(string))
		web.SessionSet(r, w, "user_nickname", userinfo["nickname"].(string))
		web.SessionSet(r, w, "user_email", userinfo["email"].(string))
		web.SessionSet(r, w, "user_id", userinfo["user_id"].(string))

		http.Redirect(w, r, r.Referer(), http.StatusFound)
		return nil
	}
	auth0.LogoutCallback = func(w http.ResponseWriter, r *http.Request) error {
		session, _ := web.Sessions.Get(r, "session")
		delete(session.Values, "_auth_id")
		delete(session.Values, "_access_token")
		delete(session.Values, "user_name")
		delete(session.Values, "user_nickname")
		delete(session.Values, "user_email")
		delete(session.Values, "user_id")
		session.Options.MaxAge = -1
		return session.Save(r, w)
	}

	stripe.Key = com.GetString("stripe_secret_key")

	return nil
}
