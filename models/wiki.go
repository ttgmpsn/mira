package models

// Wiki defines a Wiki Page
type Wiki struct {
	ContentMD    string        `json:"content_md"`
	MayRevise    bool          `json:"may_revise"`
	Reason       string        `json:"reason"`
	RevisionDate int           `json:"revision_date"`
	RevisionBy   RedditElement `json:"revision_by"`
	RevisionID   string        `json:"revision_id"`
	ContentHTML  string        `json:"content_html"`
}
