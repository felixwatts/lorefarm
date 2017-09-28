package lorefarm

import (
	"errors"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"

	"github.com/gorilla/sessions"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	DB          BookDatabase
	OAuthConfig *oauth2.Config
	SessionStore sessions.Store
)

func init() {
	var err error

	StorageBucketName = "lorefarm-181215"
	StorageBucket, err = configureStorage(StorageBucketName)

	if err != nil {
		log.Fatal(err)
	}

	OAuthConfig = configureOAuthClient("173856195020-suj3gkujjiddij6nkr77gb1jvhbdtu2m.apps.googleusercontent.com", "qntKGX2VBQzaRo0D7TBbUNtU")

	cookieStore := sessions.NewCookieStore([]byte("something-very-secret"))
	cookieStore.Options = &sessions.Options{
		HttpOnly: true,
	}
	SessionStore = cookieStore
}

func configureStorage(bucketID string) (*storage.BucketHandle, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.Bucket(bucketID), nil
}

func configureOAuthClient(clientID, clientSecret string) *oauth2.Config {
	redirectURL := os.Getenv("OAUTH2_CALLBACK")
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/oauth2callback"
	}
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
}
