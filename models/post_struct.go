package models

import "encoding/json"

// Post defines a reddit post (t3_XXXXX)
type Post struct {
	ApprovedAtUTC              int64               `json:"approved_at_utc"`
	Subreddit                  string              `json:"subreddit"`
	Selftext                   string              `json:"selftext"`
	AuthorFullname             RedditID            `json:"author_fullname"`
	Saved                      bool                `json:"saved"`
	ModReasonTitle             string              `json:"mod_reason_title"`
	Gilded                     uint                `json:"gilded"`
	Clicked                    bool                `json:"clicked"`
	Title                      string              `json:"title"`
	LinkFlairRichtext          []map[string]string `json:"link_flair_richtext"`
	SubredditNamePrefixed      string              `json:"subreddit_name_prefixed"`
	Hidden                     bool                `json:"hidden"`
	LinkFlairCSSClass          string              `json:"link_flair_css_class"`
	Downs                      int                 `json:"downs"`
	ThumbnailHeight            uint                `json:"thumbnail_height"`
	HideScore                  bool                `json:"hide_score"`
	Name                       RedditID            `json:"name"`
	Quarantine                 bool                `json:"quarantine"`
	LinkFlairTextColor         string              `json:"link_flair_text_color"`
	AuthorFlairBackgroundColor string              `json:"author_flair_background_color"`
	SubredditType              string              `json:"subreddit_type"`
	Ups                        int                 `json:"ups"`
	TotalAwardsReceived        uint                `json:"total_awards_received"`
	MediaEmbed                 json.RawMessage     `json:"media_embed"` // is: []string
	ThumbnailWidth             uint                `json:"thumbnail_width"`
	AuthorFlairTemplateID      string              `json:"author_flair_template_id"`
	IsOriginalContent          bool                `json:"is_original_content"`
	UserReports                []UserReport        `json:"user_reports"`
	SecureMedia                string              `json:"secure_media"`
	IsRedditMediaDomain        bool                `json:"is_reddit_media_domain"`
	IsMeta                     bool                `json:"is_meta"`
	Category                   string              `json:"category"`
	SecureMediaEmbed           json.RawMessage     `json:"secure_media_embed"` // is: []string
	LinkFlairText              string              `json:"link_flair_text"`
	CanModPost                 bool                `json:"can_mod_post"`
	Score                      int                 `json:"score"`
	ApprovedBy                 string              `json:"approved_by"`
	Thumbnail                  string              `json:"thumbnail"`
	Edited                     json.RawMessage     `json:"edited"` // bool or float64
	AuthorFlairCSSClass        string              `json:"author_flair_css_class"`
	AuthorFlairRichtext        []map[string]string `json:"author_flair_richtext"`
	Gildings                   map[string]int      `json:"gildings"`
	PostHint                   string              `json:"post_hint"`
	ContentCategories          []string            `json:"content_categories"`
	IsSelf                     bool                `json:"is_self"`
	ModNote                    string              `json:"mod_note"`
	Created                    float64             `json:"created"`
	LinkFlairType              string              `json:"link_flair_type"`
	BannedBy                   string              `json:"banned_by"`
	AuthorFlairType            string              `json:"author_flair_type"`
	Domain                     string              `json:"domain"`
	SelftextHTML               string              `json:"selftext_html"`
	Likes                      bool                `json:"likes"`
	SuggestedSort              string              `json:"suggested_sort"`
	BannedAtUTC                float64             `json:"banned_at_utc"`
	ViewCount                  uint                `json:"view_count"`
	Archived                   bool                `json:"archived"`
	NoFollow                   bool                `json:"no_follow"`
	IsCrosspostable            bool                `json:"is_crosspostable"`
	Pinned                     bool                `json:"pinned"`
	Over18                     bool                `json:"over_18"`
	Preview                    struct {
		Images []struct {
			Source      Image   `json:"source"`
			Resolutions []Image `json:"resolutions"`
			ID          string  `json:"id"`
		} `json:"images"`
		Enabled bool `json:"enabled"`
	} `json:"preview"`
	Awardings                []PostAward `json:"all_awardings"`
	MediaOnly                bool        `json:"media_only"`
	CanGild                  bool        `json:"can_gild"`
	Spoiler                  bool        `json:"spoiler"`
	Locked                   bool        `json:"locked"`
	AuthorFlairText          string      `json:"author_flair_text"`
	Visited                  bool        `json:"visited"`
	NumReports               int         `json:"num_reports"`
	Distinguished            bool        `json:"distinguished"`
	SubredditID              RedditID    `json:"subreddit_id"`
	ModReasonBy              string      `json:"mod_reason_by"`
	RemovalReason            string      `json:"removal_reason"`
	LinkFlairBackgroundColor string      `json:"link_flair_background_color"`
	ID                       string      `json:"id"`
	IsRobotIndexable         bool        `json:"is_robot_indexable"`
	Author                   string      `json:"author"`
	NumCrossposts            uint        `json:"num_crossposts"`
	NumComments              uint        `json:"num_comments"`
	SendReplies              bool        `json:"send_replies"`
	WhitelistStatus          string      `json:"whitelist_status"`
	ContestMode              bool        `json:"contest_mode"`
	ModReports               []ModReport `json:"mod_reports"`
	AuthorPatreonFlair       bool        `json:"author_patreon_flair"`
	AuthorFlairTextColor     string      `json:"author_flair_text_color"`
	Permalink                string      `json:"permalink"`
	ParentWhitelistStatus    string      `json:"parent_whitelist_status"`
	Stickied                 bool        `json:"stickied"`
	URL                      string      `json:"url"`
	SubredditSubscribers     uint        `json:"subreddit_subscribers"`
	CreatedUTC               float64     `json:"created_utc"`
	Media                    []string    `json:"media"`
	IsVideo                  bool        `json:"is_video"`
	Approved                 bool        `json:"approved"`
	Removed                  bool        `json:"removed"`
}

// PostAward defines a post "hipster gilding"
type PostAward struct {
	IsEnabled           bool     `json:"is_enabled"`
	Count               uint     `json:"count"`
	SubredditID         RedditID `json:"subreddit_id"`
	Description         string   `json:"description"`
	CoinReward          uint     `json:"coin_reward"`
	IconWidth           uint     `json:"icon_width"`
	IconURL             string   `json:"icon_url"`
	DaysOfPremium       uint     `json:"days_of_premium"`
	IconHeight          float64  `json:"icon_height"`
	ResizedIcons        []Image  `json:"resized_icons"`
	DaysOfDripExtension float64  `json:"days_of_drip_extension"`
	AwardType           string   `json:"award_type"`
	CoinPrice           uint     `json:"coin_price"`
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
}

// Image defines basic image properties
type Image struct {
	URL    string `json:"url"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
}
