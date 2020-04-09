package models

import (
	"fmt"
	"time"
)

// GetID returns the RedditID of the ModAction - which isn't actually a RedditID :(
func (m ModAction) GetID() RedditID { return RedditID(m.ID) }

// CreatedAt returns your account creation date
func (m ModAction) CreatedAt() time.Time { return time.Unix(m.CreatedUTC, 0) }

// GetURL returns a link to your profile
func (m ModAction) GetURL() string {
	return fmt.Sprintf("https://www.reddit.com/%s", m.TargetPermalink)
}
