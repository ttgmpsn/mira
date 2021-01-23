package models

// Stylesheet has details abouta a subreddits stylesheet.
type Stylesheet struct {
	Images []struct {
		URL  string `json:"url"`
		Link string `json:"link"`
		Name string `json:"name"`
	} `json:"images"`
	SubredditID RedditID `json:"subreddit_id"`
	Stylesheet  string   `json:"stylesheet"`
}
