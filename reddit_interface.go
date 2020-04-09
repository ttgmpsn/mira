package mira

// Interface can be used for any reddit object (Post, Comment, etc.) to avoid type hinting.
type Interface interface {
	GetID() string
	GetParentID() string
	GetTitle() string
	GetBody() string
	GetAuthor() string
	GetName() string
	GetKarma() float64
	GetUps() float64
	GetDowns() float64
	GetSubreddit() string
	GetCreated() float64
	GetFlair() string
	GetURL() string
	IsRoot() bool
}
