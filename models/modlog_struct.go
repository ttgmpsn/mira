package models

// ModAction has details about stuff a moderator bit. Usually bad stuff.
type ModAction struct {
	Description           string   `json:"description"`
	TargetBody            string   `json:"target_body"`
	ModID36               string   `json:"mod_id36"`
	CreatedUTC            float64  `json:"created_utc"`
	Subreddit             string   `json:"subreddit"`
	TargetTitle           string   `json:"target_title"`
	TargetPermalink       string   `json:"target_permalink"`
	SubredditNamePrefixed string   `json:"subreddit_name_prefixed"`
	Details               string   `json:"details"`
	Action                string   `json:"action"`
	TargetAuthor          string   `json:"target_author"`
	TargetFullname        RedditID `json:"target_fullname"`
	SrID36                string   `json:"sr_id36"`
	ID                    string   `json:"id"`
	Mod                   string   `json:"mod"`
}
