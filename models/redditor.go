package models

import (
	"fmt"
	"net/url"
	"time"
)

// GetID returns the RedditID of You
func (r Redditor) GetID() RedditID { return RedditID("t2_" + r.ID) }

// CreatedAt returns your account creation date
func (r Redditor) CreatedAt() time.Time { return time.Unix(int64(r.CreatedUTC), 0) }

// GetURL returns a link to your profile
func (r Redditor) GetURL() string {
	return fmt.Sprintf("https://www.reddit.com/u/%s", url.PathEscape(r.Name))
}
