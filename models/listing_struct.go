package models

import (
	"encoding/json"
	"fmt"
)

// Listing contains a selection of Elements.
// It also provides pagination via After & Before
type Listing struct {
	Modhash  string          `json:"modhash"`
	Dist     int             `json:"dist"`
	Children []RedditElement `json:"children"`
	After    string          `json:"after"`
	Before   string          `json:"before"`
}

// RedditElement contains a Reddit Element. The returned Element is defined by Kind.
type RedditElement struct {
	Kind RedditKind  `json:"kind"`
	Data RedditThing `json:"data"`
}

// UnmarshalJSON helps to convert RedditElement.Data into a useful object depending on Kind
func (r *RedditElement) UnmarshalJSON(data []byte) error {
	var m struct {
		Kind RedditKind      `json:"kind"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	switch m.Kind {
	case KComment:
		r.Data = &Comment{}
	case KRedditor:
		r.Data = &Redditor{}
	case KPost:
		r.Data = &Post{}
	//case KMessage:
	// not implemented
	case KSubreddit:
		r.Data = &Subreddit{}
	//case KAward:
	// not implemented
	case KModAction:
		r.Data = &ModAction{}
	default:
		return fmt.Errorf("%q is an invalid RedditKind", m.Kind)
	}

	r.Kind = m.Kind
	return json.Unmarshal(m.Data, r.Data)
}

// CommentActionResponse is returned by reddit when you create a comment (new or reply)
type CommentActionResponse struct {
	JSON struct {
		Errors []string `json:"errors"`
		Data   struct {
			Things []RedditElement `json:"things"`
		}
	} `json:"json"`
}

// PostActionResponse is returned by reddit when you create a post
type PostActionResponse struct {
	JSON struct {
		Errors []string `json:"errors"`
		Data   struct {
			Name RedditID `json:"name"`
			URL  string   `json:"url"`
		}
	} `json:"json"`
}
