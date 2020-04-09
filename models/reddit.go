package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// RedditID is a Reddit Thing ID (tX_XXXXX)
type RedditID string

// RedditKind defines the type of of a Reddit ID (tX)
type RedditKind string

// List of all possible RedditKinds
const (
	KComment   RedditKind = "t1"
	KRedditor  RedditKind = "t2"
	KPost      RedditKind = "t3"
	KMessage   RedditKind = "t4"
	KSubreddit RedditKind = "t5"
	KAward     RedditKind = "t6"
	KModAction RedditKind = "modaction"
	KUnknown   RedditKind = "tX"
)

// Type returns the RedditKind to a RedditID
func (id RedditID) Type() RedditKind {
	switch id[:3] {
	case "t1_":
		return KComment
	case "t2_":
		return KRedditor
	case "t3_":
		return KPost
	case "t4_":
		return KMessage
	case "t5_":
		return KSubreddit
	case "t6_":
		return KAward
	case "modaction":
		return KModAction
	default:
		return KUnknown
	}
}

type responseType string

const (
	rListing responseType = "Listing"
	rWiki    responseType = "wikipage"
)

// Response is a reply from the reddit API.
// It can be of different Kinds (Respond Types), and contains different Data
// (depending on Kind)
type Response struct {
	Kind responseType
	Data interface{}
}

// UnmarshalJSON helps to convert Response.Data into a useful object depending on Kind
func (r *Response) UnmarshalJSON(data []byte) error {
	var m struct {
		Kind responseType    `json:"kind"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	switch m.Kind {
	case rListing:
		r.Data = &Listing{}
	case rWiki:
		r.Data = &Wiki{}
	default:
		return fmt.Errorf("%q is an invalid ResponseType", m.Kind)
	}

	r.Kind = m.Kind
	return json.Unmarshal(m.Data, r.Data)
}

// RedditThing provides an interface to access common attributes
// of reddit stuff without having to type hint it.
type RedditThing interface {
	GetID() RedditID
	CreatedAt() time.Time
	GetURL() string
}

// Submission provides an interface to access common attributes
// of posts & comments without having to type hint it.
type Submission interface {
	RedditThing
	GetAuthor() string
	GetAuthorID() RedditID
	GetSubreddit() string
	GetSubredditID() RedditID
	GetParentID() RedditID
	GetTitle() string
	GetBody() string
	GetScore() int
	IsSticky() bool
	// Mod Stuff
	IsRemoved() bool
	IsApproved() bool
	GetBanned() SubModAction
	GetApproved() SubModAction
	GetReports() AllReports
}
