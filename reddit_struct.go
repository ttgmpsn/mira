package mira

import (
	"net/http"
	"time"
)

// Reddit holds the connection to the API for a user
type Reddit struct {
	Token    string  `json:"access_token"`
	Duration float64 `json:"expires_in"`
	Creds    Credentials
	Chain    []chainVals
	Stream   Streaming
	Values   RedditVals
	Client   *http.Client
}

// Streaming holds information
type Streaming struct {
	CommentListInterval time.Duration
	PostListInterval    time.Duration
	PostListSlice       int
}

// RedditVals holds configuration values for
type RedditVals struct {
	GetSubmissionFromCommentTries int
}

type chainVals struct {
	Name string
	Type string
}
