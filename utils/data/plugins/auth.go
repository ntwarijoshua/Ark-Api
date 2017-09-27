package plugins

import (
	"net/http"
	"encoding/base64"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"encoding/json"
	"strings"
)

var defaultRealm = "Authorization Required"

// Basic is the http basic auth
func Basic(username string, password string) beego.FilterFunc {
	secrets := func(user, pass string) bool {
		return user == username && pass == password
	}
	return NewBasicAuthenticator(secrets, defaultRealm)
}

// SecretProvider is the SecretProvider function
type SecretProvider func(user, pass string) bool

type CustomBasicAuth struct {
	Secrets SecretProvider
	Realm string

}

// CheckAuth Checks the username/password combination from the request. Returns
// either an empty string (authentication failed) or the name of the
// authenticated user.
// Supports MD5 and SHA1 password entries
func (a *CustomBasicAuth) CheckAuth(r *http.Request) string {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 || s[0] != "Basic" {
		return ""
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return ""
	}
	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return ""
	}

	if a.Secrets(pair[0], pair[1]) {
		return pair[0]
	}
	return ""
}


func (a *CustomBasicAuth) RequireAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="`+a.Realm+`"`)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(401)
	json.NewEncoder(w).Encode(map[string]string{"Error":"Invalid Credentials"})

}

func NewBasicAuthenticator(secrets SecretProvider, Realm string) beego.FilterFunc {
	return func(ctx *context.Context) {
		a := &CustomBasicAuth{Secrets: secrets, Realm: Realm}
		if username := a.CheckAuth(ctx.Request); username == "" {
			a.RequireAuth(ctx.ResponseWriter, ctx.Request)
		}
	}
}
