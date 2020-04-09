package models

// Redditor represents another user
type Redditor struct {
	IsEmployee        bool          `json:"is_employee"`
	IconImg           string        `json:"icon_img"`
	PrefShowSnoovatar bool          `json:"pref_show_snoovatar"`
	Name              string        `json:"name"`
	IsFriend          bool          `json:"is_friend"`
	Created           int64         `json:"created"`
	HasSubscribed     bool          `json:"has_subscribed"`
	HideFromRobots    bool          `json:"hide_from_robots"`
	CreatedUTC        int64         `json:"created_utc"`
	LinkKarma         float64       `json:"link_karma"`
	CommentKarma      float64       `json:"comment_karma"`
	IsGold            bool          `json:"is_gold"`
	IsMod             bool          `json:"is_mod"`
	Verified          bool          `json:"verified"`
	Subreddit         userSubreddit `json:"subreddit"`
	HasVerifiedEmail  bool          `json:"has_verified_email"`
	ID                string        `json:"id"`
}

type userSubreddit struct {
	DefaultSet                 bool      `json:"default_set"`
	UserIsContributor          bool      `json:"user_is_contributor"`
	BannerImg                  string    `json:"banner_img"`
	DisableContributorRequests bool      `json:"disable_contributor_requests"`
	UserIsBanned               bool      `json:"user_is_banned"`
	FreeFormReports            bool      `json:"free_form_reports"`
	CommunityIcon              string    `json:"community_icon"`
	ShowMedia                  bool      `json:"show_media"`
	IconColor                  string    `json:"icon_color"`
	UserIsMuted                bool      `json:"user_is_muted"`
	DisplayName                string    `json:"display_name"`
	HeaderImg                  string    `json:"header_img"` // *
	Title                      string    `json:"title"`
	Over18                     bool      `json:"over_18"`
	IconSize                   []float64 `json:"icon_size"`
	PrimaryColor               string    `json:"primary_color"`
	IconImg                    string    `json:"icon_img"`
	Description                string    `json:"description"`
	HeaderSize                 string    `json:"header_size"` // *
	RestrictPosting            bool      `json:"restrict_posting"`
	RestrictCommenting         bool      `json:"restrict_commenting"`
	Subscribers                int       `json:"subscribers"`
	IsDefaultIcon              bool      `json:"is_default_icon"`
	LinkFlairPosition          string    `json:"link_flair_position"`
	DisplayNamePrefixed        string    `json:"display_name_prefixed"`
	KeyColor                   string    `json:"key_color"`
	Name                       RedditID  `json:"name"`
	IsDefaultBanner            bool      `json:"is_default_banner"`
	URL                        string    `json:"url"`
	BannerSize                 []int     `json:"banner_size"`
	UserIsModerator            bool      `json:"user_is_moderator"`
	PublicDescription          string    `json:"public_description"`
	LinkFlairEnabled           bool      `json:"link_flair_enabled"`
	SubredditType              string    `json:"subreddit_type"`
	UserIsSubscriber           bool      `json:"user_is_subscriber"`
}
