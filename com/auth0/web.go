package auth0

import (
	"fmt"
	"net/http"
	"net/url"
	"text/template"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/gliderlabs.io/com/web"
)

func (c *Component) WebTemplateFuncMap(r *http.Request) template.FuncMap {
	return template.FuncMap{
		"auth0": func() string {
			return fmt.Sprintf(`
				var auth0 = {};
				auth0.init = function() {
					var js = document.createElement("script");
					js.type = "text/javascript";
					js.src = "https://cdn.auth0.com/js/lock-9.0.min.js";
					js.onload = function() {
						auth0.lock = new Auth0Lock('%s', '%s');
						auth0.login = function() {
							auth0.lock.show({callbackURL: '%s'});
						};
					};
					document.body.appendChild(js);
				};
				auth0.init();
			`,
				com.GetString("client_id"),
				com.GetString("domain"),
				com.GetString("callback_url"))
		},
	}
}

func (c *Component) MatchHTTP(r *http.Request) bool {
	if cb, err := url.Parse(com.GetString("callback_url")); err == nil {
		if r.URL.Path == cb.Path {
			return true
		}
	}
	if logout, err := url.Parse(com.GetString("logout_url")); err == nil {
		if r.URL.Path == logout.Path {
			return true
		}
	}
	return false
}

// ServeHTTP of web.Handler extension point
func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if logout, err := url.Parse(com.GetString("logout_url")); err == nil {
		if r.URL.Path == logout.Path {
			if web.Return(r, w) {
				return
			}
			if err := web.SetReturn(r, w, r.Referer()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if LogoutCallback != nil {
				if err := LogoutCallback(w, r); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			http.Redirect(w, r, Client().LogoutURL(com.GetString("logout_url")), http.StatusFound)
			return
		}
	}

	token, err := Client().NewToken(r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if LoginCallback != nil {
		if err := LoginCallback(w, r, token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
