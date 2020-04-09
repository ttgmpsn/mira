package models

type CommentWrap struct {
	JSON CommentJSON `json:"json"`
}

type CommentJSON struct {
	Errors []string        `json:"errors"`
	Data   CommentJSONData `json:"data"`
}

type CommentJSONData struct {
	Things []CommentJSONDataThing `json:"things"`
}

type CommentJSONDataThing struct {
	Kind string                   `json:"kind"`
	Data CommentJSONDataThingData `json:"data"`
}

type CommentJSONDataThingData struct {
	AuthorFlairBackgroundColor string       `json:"author_flair_background_color"`
	TotalAwardsReceived        float64      `json:"total_awards_received"`
	ApprovedAtUTC              float64      `json:"approved_at_utc"`
	Distinguished              string       `json:"distinguished"`
	ModReasonBy                string       `json:"mod_reason_by"`
	BannedBy                   string       `json:"banned_by"`
	AuthorFlairType            string       `json:"author_flair_type"`
	RemovalReason              string       `json:"removal_reason"`
	LinkID                     string       `json:"link_id"`
	AuthorFlairTemplateID      string       `json:"author_flair_template_id"`
	Likes                      bool         `json:"likes"`
	Replies                    string       `json:"replies"`
	UserReports                []UserReport `json:"user_reports"`
	Saved                      bool         `json:"saved"`
	ID                         string       `json:"id"`
	BannedAtUTC                string       `json:"banned_at_utc"`
	ModReasonTitle             string       `json:"mod_reason_title"`
	Gilded                     float64      `json:"gilded"`
	Archived                   bool         `json:"archived"`
	NoFollow                   bool         `json:"no_follow"`
	Author                     string       `json:"author"`
	RteMode                    string       `json:"rte_mode"`
	CanModPost                 bool         `json:"can_mod_post"`
	CreatedUTC                 float64      `json:"created_utc"`
	SendReplies                bool         `json:"send_replies"`
	ParentID                   string       `json:"parent_id"`
	Score                      float64      `json:"score"`
	AuthorFullname             string       `json:"author_fullname"`
	ApprovedBy                 string       `json:"approved_by"`
	ModNote                    string       `json:"mod_note"`
	AllAwardings               []string     `json:"all_awardings"`
	SubredditID                string       `json:"subreddit_id"`
	Body                       string       `json:"body"`
	Edited                     bool         `json:"edited"`
	Gildings                   Gilding      `json:"gildings"`
	AuthorFlairCSSClass        string       `json:"author_flair_css_class"`
	Name                       string       `json:"name"`
	AuthorPatreonFlair         bool         `json:"author_patreon_flair"`
	Downs                      float64      `json:"downs"`
	AuthorFlairRichtext        []string     `json:"author_flair_richtext"`
	IsSubmitter                bool         `json:"is_submitter"`
	CollapsedReason            string       `json:"collapsed_reason"`
	BodyHTML                   string       `json:"body_html"`
	Stickied                   bool         `json:"stickied"`
	CanGild                    bool         `json:"can_gild"`
	Removed                    bool         `json:"removed"`
	Approved                   bool         `json:"approved"`
	AuthorFlairTextColor       string       `json:"author_flair_text_color"`
	ScoreHidden                bool         `json:"score_hidden"`
	Permalink                  string       `json:"permalink"`
	NumReports                 float64      `json:"num_reports"`
	Locked                     bool         `json:"locked"`
	ReportReasons              []string     `json:"report_reasons"`
	Created                    float64      `json:"created"`
	Subreddit                  string       `json:"subreddit"`
	AuthorFlairText            string       `json:"author_flair_text"`
	Spam                       bool         `json:"spam"`
	Collapsed                  bool         `json:"collapsed"`
	SubredditNamePrefixed      string       `json:"subreddit_name_prefixed"`
	Controversiality           float64      `json:"controversiality"`
	IgnoreReports              bool         `json:"ignore_reports"`
	ModReports                 []ModReport  `json:"mod_reports"`
	SubredditType              string       `json:"subreddit_type"`
	Ups                        float64      `json:"ups"`
}

type Gilding struct {
	Gid map[string]int `json:"gid"`
}
