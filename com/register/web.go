package register

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/gliderlabs/gliderlabs.io/com/auth0"
	"github.com/gliderlabs/gliderlabs.io/com/mailgun"
	"github.com/gliderlabs/gliderlabs.io/com/web"
	auth0lib "github.com/gliderlabs/pkg/auth0"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

func (c *Component) WebTemplateFuncMap(r *http.Request) template.FuncMap {
	return template.FuncMap{
		"title": strings.Title,
	}
}

func (c *Component) MatchHTTP(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, "/register")
}

func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/register" {
		http.Redirect(w, r, "/register/", http.StatusMovedPermanently)
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/register/pay", payHandler)
	mux.HandleFunc("/register/verify/", verifyHandler)
	mux.HandleFunc("/register/", registerHandler)
	mux.ServeHTTP(w, r)
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	key := parts[len(parts)-1]
	project := parts[len(parts)-2]
	lookupKey := ExtractLookupKey(key)
	user, err := LookupUser(lookupKey)
	if err != nil {
		http.Error(w, "User lookup failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	if _, ok := user.AppMetadata.Registered[project]; !ok {
		http.Error(w, "Registration not found", http.StatusNotFound)
		return
	}
	if key != user.AppMetadata.Registered[project].Key {
		http.Error(w, "Registration key does not match", http.StatusNotFound)
		return
	}
	w.Write([]byte("OK\n"))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	project := parts[len(parts)-1]
	if project == "" {
		http.Error(w, "No project to register", http.StatusNotFound)
		return
	}

	uid := web.SessionValue(r, "user_id")
	if uid != "" {
		if key := registeredKey(uid, project); key != "" {
			web.RenderTemplate(w, r, "key", map[string]interface{}{
				"Name":     strings.Split(web.SessionValue(r, "user_name"), " ")[0],
				"Nickname": web.SessionValue(r, "user_nickname"),
				"Project":  project,
				"Key":      key,
			})
			return
		}
		web.RenderTemplate(w, r, "pay", map[string]interface{}{
			"Name":      strings.Split(web.SessionValue(r, "user_name"), " ")[0],
			"Nickname":  web.SessionValue(r, "user_nickname"),
			"Email":     web.SessionValue(r, "user_email"),
			"Project":   project,
			"StripeKey": com.GetString("stripe_pub_key"),
		})
		return
	}
	web.RenderTemplate(w, r, "login", map[string]interface{}{
		"Project": project,
	})
}

func payHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if uid := web.SessionValue(r, "user_id"); uid == "" {
		http.Error(w, "No user session", http.StatusBadRequest)
		return
	}

	key := GenerateKey(web.SessionValue(r, "user_id"))

	chargeId, err := processPayment(r, key)
	if err != nil {
		http.Error(w, "payment failed: "+err.Error(), http.StatusPaymentRequired)
		return
	}

	err = mailgun.SendText(
		r.FormValue("email"),
		"Jeff Lindsay <jeff@gliderlabs.com>",
		fmt.Sprintf("Registration Key for %s", strings.Title(r.FormValue("project"))),
		fmt.Sprintf(`Hello %s!

Thanks again for registering %s and supporting our work. Here is your registration key:

%s

For further instructions see %s

-jeff`, web.SessionValue(r, "user_name"),
			strings.Title(r.FormValue("project")),
			key,
			"http://localhost:8080/register/"+r.FormValue("project")),
	)
	if err != nil {
		log.Info("mailgun", err)
	}

	err = auth0.Client().PatchUser(web.SessionValue(r, "user_id"), auth0lib.User{
		"app_metadata": map[string]interface{}{
			"registered": map[string]interface{}{
				r.FormValue("project"): map[string]interface{}{
					"key":    key,
					"charge": chargeId,
				},
			},
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/register/"+r.FormValue("project"), http.StatusFound)
}

func registeredKey(uid, project string) string {
	user, err := LookupUser(uid)
	if err != nil {
		log.Info("auth0", err)
		return ""
	}
	if _, ok := user.AppMetadata.Registered[project]; !ok {
		return ""
	}
	return user.AppMetadata.Registered[project].Key
}

func processPayment(r *http.Request, key string) (string, error) {
	amount, err := strconv.Atoi(strings.Trim(r.FormValue("price"), "$ "))
	if err != nil {
		return "", err
	}
	chargeParams := &stripe.ChargeParams{
		Amount:   uint64(amount * 100),
		Currency: "usd",
		Desc:     "Registration of " + strings.Title(r.FormValue("project")),
		Email:    r.FormValue("email"),
	}
	chargeParams.AddMeta("key", key)
	chargeParams.AddMeta("project", r.FormValue("project"))
	chargeParams.AddMeta("charity", r.FormValue("charity"))
	chargeParams.SetSource(r.FormValue("stripeToken"))
	ch, err := charge.New(chargeParams)
	if err != nil {
		return "", err
	}
	return ch.ID, nil
}
