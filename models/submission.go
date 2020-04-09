package models

func (s *Submission) GetID() string           { return s.JSON.Data.Name }
func (s *Submission) GetDraftsCount() float64 { return s.JSON.Data.DraftsCount }
func (s *Submission) GetURL() string          { return s.JSON.Data.URL }
