package models

import "time"

// GetID returns the Subreddit RedditID
func (s Subreddit) GetID() RedditID { return s.Name }

// GetURL returns the Subreddit URL
func (s Subreddit) GetURL() string { return s.URL }

// CreatedAt returns the time.Time of the Subreddits creation
func (s Subreddit) CreatedAt() time.Time { return time.Unix(s.CreatedUTC, 0) }
