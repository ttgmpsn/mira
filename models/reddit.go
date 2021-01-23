package models

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// RedditID is a Reddit Thing ID (tX_XXXXX)
type RedditID string

// UnmarshalJSON defines an Unmarshaller for RedditID
// This is mainly used to convert int values (used in New Modmail responses) to the base36 one used elsewhere.
func (r *RedditID) UnmarshalJSON(data []byte) error {
	var t interface{}
	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}
	switch v := t.(type) {
	case string:
		*r = RedditID(v)
		return nil
	case int64:
		// Only known to be used for authors
		*r = RedditID("t2_" + big.NewInt(v).Text(36))
		return nil
	case float64:
		// Only known to be used for authors
		f := big.NewFloat(v)
		i, _ := f.Int(nil)
		*r = RedditID("t2_" + i.Text(36))
		return nil
	case nil:
		r = nil
		return nil
	}

	return fmt.Errorf("Unknown type for RedditID %v", t)
}

// Type returns the RedditKind to a RedditID
func (r RedditID) Type() RedditKind {
	s := string(r)
	switch strings.ToLower(s[:strings.IndexByte(s, '_')]) {
	case "t1":
		return KComment
	case "t2":
		return KRedditor
	case "t3":
		return KPost
	case "t4":
		return KMessage
	case "t5":
		return KSubreddit
	case "t6":
		return KAward
	case "modaction":
		return KModAction
	default:
		return KUnknown
	}
}

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

type responseType string

const (
	rListing    responseType = "Listing"
	rWiki       responseType = "wikipage"
	rStylesheet responseType = "stylesheet"
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
	//fmt.Println(string(data))
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
	case rStylesheet:
		r.Data = &Stylesheet{}
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
	GetCreated() time.Time
	// Mod Stuff
	IsRemoved() bool
	IsApproved() bool
	GetBanned() SubModAction
	GetApproved() SubModAction
	GetReports() AllReports
}
