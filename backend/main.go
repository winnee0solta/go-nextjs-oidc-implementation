package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var (
	clientID     = "clientID"
	clientSecret = "clientSecret"
	tenantID     = "tenantID"
	redirectURL  = "http://localhost:8080/oidc/login/callback"
	frontendURL  = "http://localhost:3000"
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	oauth2Config oauth2.Config
	store        = sessions.NewCookieStore([]byte("supersecurestring"))
)

type User struct {
	Email string
	Name  string
}

func init() {
	gob.Register(&User{}) // Required for storing in session
}

func main() {

	ctx := context.Background()

	providerURL := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", tenantID)
	var err error
	provider, err = oidc.NewProvider(ctx, providerURL)
	if err != nil {
		log.Fatalf("Failed to initialize provider: %v", err)
	}

	verifier = provider.Verifier(&oidc.Config{ClientID: clientID})

	oauth2Config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		RedirectURL:  redirectURL,
	}

	http.HandleFunc("/oidc/login", handleLogin)
	http.HandleFunc("/oidc/login/callback", handleCallback)
	http.HandleFunc("/oidc/logout", handleLogout)
	http.HandleFunc("/me", handleMe)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauth2Config.AuthCodeURL("state-random", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.URL.Query().Get("state") != "state-random" {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "token exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "missing id_token", http.StatusInternalServerError)
		return
	}

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "token verification failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var claims struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "failed to parse claims: "+err.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := store.Get(r, "session")
	session.Values["user"] = &User{Email: claims.Email, Name: claims.Name}
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 8,
		HttpOnly: true,
		Secure:   true, // Set to false for local testing over HTTP
		SameSite: http.SameSiteLaxMode,
	}
	session.Save(r, w)

	http.Redirect(w, r, frontendURL+"/dashboard", http.StatusFound)
}

func handleMe(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	user, ok := session.Values["user"].(*User)
	if !ok || user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)

	// Azure logout URL
	logoutURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/logout?post_logout_redirect_uri=%s", tenantID, frontendURL)
	http.Redirect(w, r, logoutURL, http.StatusFound)
}
