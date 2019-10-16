package main

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/auth"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/auth/providers/facebook"
	"github.com/qor/auth/providers/github"
	"github.com/qor/auth/providers/google"
	"github.com/qor/auth/providers/password"
	"github.com/qor/auth/providers/twitter"
	"github.com/qor/session/manager"

	_ "github.com/mattn/go-sqlite3"

	"net/http"
)

var (
	// Initialize gorm DB
	gormDB, _ = gorm.Open("sqlite3", "sample.db")

	// Auth Initialize Auth with configuration
	Auth = auth.New(&auth.Config{
		DB: gormDB,
	})
)

func init() {
	// Migrate AuthIdentity model, AuthIdentity will be used to save auth info, like username/password, oauth token, you could change that.
	gormDB.AutoMigrate(&auth_identity.AuthIdentity{})

	// Register Auth providers
	// Allow use username/password
	Auth.RegisterProvider(password.New(&password.Config{}))

	// Allow use Github
	Auth.RegisterProvider(github.New(&github.Config{
		ClientID:     "github client id",
		ClientSecret: "github client secret",
	}))

	// Allow use Google
	Auth.RegisterProvider(google.New(&google.Config{
		ClientID:       "google client id",
		ClientSecret:   "google client secret",
		AllowedDomains: []string{}, // Accept all domains, instead you can pass a whitelist of acceptable domains
	}))

	// Allow use Facebook
	Auth.RegisterProvider(facebook.New(&facebook.Config{
		ClientID:     "facebook client id",
		ClientSecret: "facebook client secret",
	}))

	// Allow use Twitter
	Auth.RegisterProvider(twitter.New(&twitter.Config{
		ClientID:     "twitter client id",
		ClientSecret: "twitter client secret",
	}))
}

func main() {
	mux := http.NewServeMux()

	// Mount Auth to Router
	mux.Handle("/auth/", Auth.NewServeMux())
	http.ListenAndServe(":9000", manager.SessionManager.Middleware(mux))
}
