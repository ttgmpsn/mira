package models

import "encoding/json"

// Comment defines a reddit comment (t1_XXXXX)
type Comment struct {
	AuthorFlairBackgroundColor string              `json:"author_flair_background_color"`
	TotalAwardsReceived        uint                `json:"total_awards_received"`
	ApprovedAtUTC              float64             `json:"approved_at_utc"`
	Distinguished              string              `json:"distinguished"`
	ModReasonBy                string              `json:"mod_reason_by"`
	BannedBy                   json.RawMessage     `json:"banned_by"`
	AuthorFlairType            string              `json:"author_flair_type"`
	RemovalReason              string              `json:"removal_reason"`
	LinkID                     RedditID            `json:"link_id"`
	AuthorFlairTemplateID      string              `json:"author_flair_template_id"`
	Likes                      bool                `json:"likes"`
	Replies                    json.RawMessage     `json:"replies"` // is []Response, but reddit likes to fuck with us :(
	UserReports                []UserReport        `json:"user_reports"`
	Saved                      bool                `json:"saved"`
	ID                         string              `json:"id"`
	BannedAtUTC                float64             `json:"banned_at_utc"`
	ModReasonTitle             string              `json:"mod_reason_title"`
	Gilded                     uint                `json:"gilded"`
	Archived                   bool                `json:"archived"`
	NoFollow                   bool                `json:"no_follow"`
	Author                     string              `json:"author"`
	RteMode                    string              `json:"rte_mode"`
	CanModPost                 bool                `json:"can_mod_post"`
	CreatedUTC                 float64             `json:"created_utc"`
	SendReplies                bool                `json:"send_replies"`
	ParentID                   RedditID            `json:"parent_id"`
	Score                      int                 `json:"score"`
	AuthorFullname             RedditID            `json:"author_fullname"`
	ApprovedBy                 string              `json:"approved_by"`
	ModNote                    string              `json:"mod_note"`
	AllAwardings               []PostAward         `json:"all_awardings"`
	SubredditID                RedditID            `json:"subreddit_id"`
	Body                       string              `json:"body"`
	Edited                     json.RawMessage     `json:"edited"` // bool or float64
	Gildings                   map[string]int      `json:"gildings"`
	AuthorFlairCSSClass        string              `json:"author_flair_css_class"`
	Name                       RedditID            `json:"name"`
	AuthorPatreonFlair         bool                `json:"author_patreon_flair"`
	Downs                      int                 `json:"downs"`
	AuthorFlairRichtext        []map[string]string `json:"author_flair_richtext"`
	IsSubmitter                bool                `json:"is_submitter"`
	CollapsedReason            string              `json:"collapsed_reason"`
	BodyHTML                   string              `json:"body_html"`
	Stickied                   bool                `json:"stickied"`
	CanGild                    bool                `json:"can_gild"`
	Removed                    bool                `json:"removed"`
	Approved                   bool                `json:"approved"`
	AuthorFlairTextColor       string              `json:"author_flair_text_color"`
	ScoreHidden                bool                `json:"score_hidden"`
	Permalink                  string              `json:"permalink"`
	NumReports                 int                 `json:"num_reports"`
	Locked                     bool                `json:"locked"`
	Created                    float64             `json:"created"`
	Subreddit                  string              `json:"subreddit"`
	AuthorFlairText            string              `json:"author_flair_text"`
	Spam                       bool                `json:"spam"`
	Collapsed                  bool                `json:"collapsed"`
	SubredditNamePrefixed      string              `json:"subreddit_name_prefixed"`
	Controversiality           int                 `json:"controversiality"`
	IgnoreReports              bool                `json:"ignore_reports"`
	ModReports                 []ModReport         `json:"mod_reports"`
	SubredditType              string              `json:"subreddit_type"`
	Ups                        int                 `json:"ups"`
}
