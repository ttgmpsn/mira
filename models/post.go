package models

func (p *PostListing) GetChildren() []PostListingChild  { return p.Data.Children }
func (plc PostListingChild) GetSubreddit() string       { return plc.Data.Subreddit }
func (plc PostListingChild) GetSubredditID() string     { return plc.Data.SubredditID }
func (plc PostListingChild) GetName() string            { return plc.Data.Title }
func (plc PostListingChild) GetTitle() string           { return plc.Data.Title }
func (plc PostListingChild) GetID() string              { return plc.Data.Name }
func (plc PostListingChild) GetParentID() string        { return plc.Data.Name }
func (plc PostListingChild) GetAuthor() string          { return plc.Data.Author }
func (plc PostListingChild) GetAuthorID() string        { return plc.Data.AuthorFullname }
func (plc PostListingChild) GetCreated() float64        { return plc.Data.CreatedUTC }
func (plc PostListingChild) GetKarma() float64          { return plc.Data.Ups - plc.Data.Downs }
func (plc PostListingChild) GetUps() float64            { return plc.Data.Ups }
func (plc PostListingChild) GetDowns() float64          { return plc.Data.Downs }
func (plc PostListingChild) GetScore() float64          { return plc.Data.Score }
func (plc PostListingChild) GetBody() string            { return plc.Data.Selftext }
func (plc PostListingChild) GetAuthorFlair() string     { return plc.Data.AuthorFlairText }
func (plc PostListingChild) GetPermalink() string       { return plc.Data.Permalink }
func (plc PostListingChild) GetURL() string             { return plc.Data.URL }
func (plc PostListingChild) GetFlair() string           { return plc.Data.LinkFlairText }
func (plc PostListingChild) GetCommentCount() float64   { return plc.Data.NumComments }
func (plc PostListingChild) GetCrosspostCount() float64 { return plc.Data.NumCrossposts }
func (plc PostListingChild) IsRoot() bool               { return true }
