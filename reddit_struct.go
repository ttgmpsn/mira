package mira

import (
	"context"
	"net/http"
	"time"

	"github.com/ttgmpsn/mira/models"
	"golang.org/x/oauth2"
)

// Credentials stores information for authing towards the Reddit API
type Credentials struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
	UserAgent    string
	RedirectURL  string
}

// Reddit holds the connection to the API for a user
type Reddit struct {
	Client      *http.Client
	creds       Credentials
	OAuthConfig *oauth2.Config
	TokenExpiry time.Time
	UserAgent   string
	ctx         context.Context

	Chain  []chainVals
	Values redditVals
}

type redditVals struct {
	GetSubmissionFromCommentTries int
}

type chainVals struct {
	Name string
	Type models.RedditKind
}
