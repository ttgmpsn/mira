package models

import "time"

// GetID returns the RedditID of the Post
func (p Post) GetID() RedditID { return p.Name }

// GetSubreddit returns the name of the Subreddit the post was posted in
func (p Post) GetSubreddit() string { return p.Subreddit }

// GetSubredditID returns the RedditID of the subreddit the post was posted in
func (p Post) GetSubredditID() RedditID { return p.SubredditID }

// GetParentID returns the parent RedditID of the Post
func (p Post) GetParentID() RedditID { return p.GetID() }

// GetAuthor returns the name of the Post Author
func (p Post) GetAuthor() string { return p.Author }

// GetAuthorID returns the RedditID of the Post Author
func (p Post) GetAuthorID() RedditID { return p.AuthorFullname }

// CreatedAt returns time.Time the Post was created at
func (p Post) CreatedAt() time.Time { return time.Unix(p.CreatedUTC, 0) }

// GetBody returns the content of the Post in Markdown
func (p Post) GetBody() string { return p.Selftext }

// GetScore returns the current score of the Post
func (p Post) GetScore() int { return p.Score }

// IsSticky tells you if the Post has been stickied
func (p Post) IsSticky() bool { return p.Stickied }

// IsRemoved tells you if the post has been removed (only if you are mod)
func (p Post) IsRemoved() bool { return p.Removed }

// IsApproved tells you if the post has been approved (only if you are mod)
func (p Post) IsApproved() bool { return p.Approved }

// GetTitle returns the title of the Post
func (p Post) GetTitle() string { return p.Title }

// GetURL returns the link to the Post
func (p Post) GetURL() string { return p.URL }

// GetBannedBy returns the user who deleted the Post
func (p Post) GetBannedBy() string { return p.BannedBy }

// GetNumReports returns the report count for the Post
func (p Post) GetNumReports() uint { return p.NumReports }
