package models

import (
	"fmt"
	"net/url"
	"time"
)

// GetID returns the RedditID of You
func (me Me) GetID() RedditID { return RedditID("t2_" + me.ID) }

// CreatedAt returns your account creation date
func (me Me) CreatedAt() time.Time { return time.Unix(int64(me.CreatedUTC), 0) }

// GetURL returns a link to your profile
func (me Me) GetURL() string {
	return fmt.Sprintf("https://www.reddit.com/u/%s", url.PathEscape(me.Name))
}
