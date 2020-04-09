package models

type Wiki struct {
	Kind string   `json:"kind"`
	Data WikiData `json:"data"`
}

type WikiData struct {
	ContentMD    string   `json:"content_md"`
	MayRevise    bool     `json:"may_revise"`
	Reason       string   `json:"reason"`
	RevisionDate int      `json:"revision_date"`
	RevisionBy   Redditor `json:"revision_by"`
	RevisionID   string   `json:"revision_id"`
	ContentHTML  string   `json:"content_html"`
}
