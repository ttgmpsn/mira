package mira

import (
	"context"
	"net/http"
	"time"

	"github.com/ttgmpsn/mira/models"
	"golang.org/x/oauth2"
)

// Credentials stores information for authing towards the Reddit API.
// You can create an app here: https://old.reddit.com/prefs/apps/
type Credentials struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
	UserAgent    string
	RedirectURL  string
}

// Reddit holds the connection to the API for a user. You can have multiple Reddit instances at the same time (see Example below).
//
// Calling Methods
//
// Each method is called similarly: First, you queue an object (Subreddit, Comment, Redditor etc) you want to perform the action on.
// Afterwards (this can be done in the same call), you select the command:
//  reddit.Subreddit("iama").Submit("I just did a bot, AMA", "Hey all! I just created a bot. AMA.")
//  reddit.Redditor("spez").Compose("Hi spez!", "Hi spez how are you?")
//
// Working with submissions
//
// If you call a function that returns a post or comment, you will receive an object that confirms to the models.Submission interface.
// You should be able to get most of the values without type hinting with the functions provided in the interface:
//  submission = reddit.SubmissionInfoID("t3_aaaaa")
//  fmt.Println("Got a post with body:", post.GetBody())
//
// If you need detailed information, you can do a type assertion:
//  post, _ := submission.(*miramodels.Post)
//  fmt.Println("Post Flair Text:", post.LinkFlairText)
type Reddit struct {
	Client      *http.Client
	creds       Credentials
	OAuthConfig *oauth2.Config
	TokenExpiry time.Time
	UserAgent   string
	ctx         context.Context

	Chain  chan *chainVals
	Values redditVals
}

type redditVals struct {
	GetSubmissionFromCommentTries int
}

type chainVals struct {
	Name string
	Type models.RedditKind
}
