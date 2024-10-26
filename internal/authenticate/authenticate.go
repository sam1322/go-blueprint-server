package authenticate

import (
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"log"
	"os"
)

const (
	key    = "randomString"
	maxAge = 60 * 60 * 24 * 7 // 1 week
	isProd = false
)

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecretKey := os.Getenv("GOOGLE_CLIENT_SECRET_KEY")

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	googleProvider := google.New(googleClientId, googleClientSecretKey, "http://localhost:8080/auth/google/callback")
	googleProvider.SetPrompt("consent", "select_account")
	//githubProvider := github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:8080/auth/github/callback")
	githubProvider := github.New(os.Getenv("GITHUB_OAUTH_CLIENT_ID"), os.Getenv("GITHUB_OAUTH_CLIENT_SECRET"), "http://localhost:8080/auth/github/callback", "user:email")

	goth.UseProviders(
		googleProvider,
		githubProvider,
		//google.New(googleClientId, googleClientSecretKey, "http://localhost:8080/auth/google/callback"),
		//github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:8080/auth/github/callback"),
	)

}
