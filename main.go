package main

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"strings"
)

var conf oauth2.Config
var verifier *oidc.IDTokenVerifier

func main() {
	tenantID := os.Getenv("AZURE_AD_TENANT_ID")
	clientID := os.Getenv("AZURE_AD_CLIENT_ID")

	issuer := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", tenantID)
	provider, err := oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		log.Fatal(err)
	}
	verifier = provider.Verifier(&oidc.Config{ClientID: clientID})

	conf = oauth2.Config{
		ClientID: clientID,
		ClientSecret: "you_wont_use_it",
		RedirectURL: "http://localhost:8000/redirect",
		Endpoint: provider.Endpoint(),
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", fuga)
	mux.HandleFunc("/redirect", hoge)
	server := &http.Server{
		Addr: "localhost:8000",
		Handler: mux,
	}
	_ = server.ListenAndServe()
}

func fuga(w http.ResponseWriter, r *http.Request) {
	url := conf.AuthCodeURL("random_state",
		oauth2.SetAuthURLParam("nonce", "random_nonce"),
		oauth2.SetAuthURLParam("response_type", "id_token"),
		oauth2.SetAuthURLParam("response_mode", "form_post"),
	)
	http.Redirect(w, r, url, 301)
}

func hoge(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, r.ContentLength)
	_, _ = r.Body.Read(b)

	firstParam := strings.Split(string(b), "&")[0]
	secondParam := strings.Split(string(b), "&")[1]

	if "error" == strings.Split(firstParam, "=")[0] {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintln(w, strings.Split(secondParam, "=")[1])
	}

	idToken := strings.Split(firstParam, "=")[1]
	if _, err := verifier.Verify(r.Context(), idToken); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintln(w, "invalid id token")
	}
}