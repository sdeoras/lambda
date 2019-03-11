package login

import (
	"gan/src/env"
	"gan/src/route"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/sdeoras/oauth"
)

const (
	GoogleProvider = "google"
)

var (
	once     sync.Once
	Provider map[string]oauth.Provider
)

func init() {
	once.Do(func() {
		Provider = make(map[string]oauth.Provider)
		redirectUrl := "https://" + filepath.Join(
			env.Domain,
			env.FuncName,
			route.OAuthGoogleCallback)
		Provider[GoogleProvider] = oauth.NewGoogleProvider(
			redirectUrl, env.ClientId, env.ClientSecret)
	})
}

func Google(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, Provider[GoogleProvider].Url(), http.StatusTemporaryRedirect)
}
