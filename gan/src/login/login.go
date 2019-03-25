package login

import (
	"gan/src/config"
	"gan/src/route"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/sdeoras/oauth"
)

const (
	GoogleProvider = "google"
	https          = "https://"
)

var (
	once     sync.Once
	Provider map[string]oauth.Provider
)

func init() {
	once.Do(func() {
		Provider = make(map[string]oauth.Provider)
		redirectUrl := https + filepath.Join(
			config.Config.Domain,
			config.Config.FuncName,
			route.OAuthGoogleCallback)
		Provider[GoogleProvider] = oauth.NewGoogleProvider(
			redirectUrl,
			config.Config.Oauth.Google.ClientId,
			config.Config.Oauth.Google.ClientSecret,
		)
	})
}

func Google(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, Provider[GoogleProvider].Url(), http.StatusTemporaryRedirect)
}
