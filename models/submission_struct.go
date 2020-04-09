package models

type Submission struct {
	JSON SubmissionJSON `json:"json"`
}

type SubmissionJSON struct {
	Errors []string           `json:"errors"`
	Data   SubmissionJSONData `json:"data"`
}

type SubmissionJSONData struct {
	URL         string  `json:"url"`
	DraftsCount float64 `json:"drafts_count"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
}
