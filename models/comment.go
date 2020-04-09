package models

func (c CommentWrap) getThing() CommentJSONDataThing {
	if len(c.JSON.Data.Things) > 0 {
		return c.JSON.Data.Things[0]
	}
	return CommentJSONDataThing{}
}
func (c CommentWrap) GetID() string          { return c.getThing().Data.Name }
func (c CommentWrap) GetSubredditID() string { return c.getThing().Data.SubredditID }
func (c CommentWrap) GetParentID() string    { return c.getThing().Data.ParentID }
func (c CommentWrap) GetAuthor() string      { return c.getThing().Data.Author }
func (c CommentWrap) GetAuthorID() string    { return c.getThing().Data.AuthorFullname }
func (c CommentWrap) GetSubreddit() string   { return c.getThing().Data.Subreddit }
func (c CommentWrap) CreatedAt() float64     { return c.getThing().Data.CreatedUTC }
func (c CommentWrap) GetBody() string        { return c.getThing().Data.Body }
func (c CommentWrap) GetScore() float64      { return c.getThing().Data.Score }
func (c CommentWrap) GetUps() float64        { return c.getThing().Data.Ups }
func (c CommentWrap) GetDowns() float64      { return c.getThing().Data.Downs }
func (c CommentWrap) IsSticky() bool         { return c.getThing().Data.Stickied }
func (c CommentWrap) IsRemoved() bool        { return c.getThing().Data.Removed }
func (c CommentWrap) IsApproved() bool       { return c.getThing().Data.Approved }
func (c CommentWrap) IsAuthor() bool         { return c.getThing().Data.IsSubmitter }
