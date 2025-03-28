package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"fmt"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

const (
	key = "something"
	MaxAge = 86400 * 30
	IsProd = false
)

type Response struct {
    Status  string `json:"status"`
    Message string `json:"message"`
    User    goth.User `json:"user,omitempty"`
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
    gothic.BeginAuthHandler(w, r)
}

func NewAuth() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("skill issue lol")
    }

    store := sessions.NewCookieStore([]byte(key))
    store.MaxAge(MaxAge)
    store.Options.Path = "/"
    store.Options.HttpOnly = true
    store.Options.Secure = IsProd

    gothic.Store = store

    goth.UseProviders(
        github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://127.0.0.1:8080/auth/github/callback"),
    )

    p := pat.New()

    p.Use(corsMiddleware)

	// Callback
    p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
        user, err := gothic.CompleteUserAuth(res, req)
        if err != nil {
            res.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(res).Encode(Response{Status: "error", Message: err.Error()})
            return
        }

        json.NewEncoder(res).Encode(Response{Status: "success", User: user})
    })

	// Logout
    p.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
        gothic.Logout(res, req)
        http.Redirect(res, req, os.Getenv("FRONTEND_URL")+"/logout", http.StatusTemporaryRedirect)
    })

	//Login
    p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
        if _, err := gothic.CompleteUserAuth(res, req); err == nil {
            http.Redirect(res, req, os.Getenv("FRONTEND_URL")+"/home", http.StatusTemporaryRedirect)
        } else {
            gothic.BeginAuthHandler(res, req)
        }
    })

fmt.Println("GITHUB_KEY:", os.Getenv("GITHUB_KEY"))
fmt.Println("GITHUB_SECRET:", os.Getenv("GITHUB_SECRET"))
    http.ListenAndServe(":8080", p)
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    })
}

