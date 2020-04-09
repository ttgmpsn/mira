package models

import (
	"fmt"
	"time"
)

// GetID returns the RedditID of the Comment
func (c Comment) GetID() RedditID { return c.Name }

// GetSubreddit returns the name of the Subreddit the comment was posted in
func (c Comment) GetSubreddit() string { return c.Subreddit }

// GetSubredditID returns the RedditID of the subreddit the comment was posted in
func (c Comment) GetSubredditID() RedditID { return c.SubredditID }

// GetParentID returns the parent RedditID of the Comment
func (c Comment) GetParentID() RedditID { return c.ParentID }

// GetAuthor returns the name of the Comment Author
func (c Comment) GetAuthor() string { return c.Author }

// GetAuthorID returns the RedditID of the Comment Author
func (c Comment) GetAuthorID() RedditID { return c.AuthorFullname }

// CreatedAt returns the time.Time the post was created at
func (c Comment) CreatedAt() time.Time { return time.Unix(c.CreatedUTC, 0) }

// GetBody returns the content of the Comment in Markdown
func (c Comment) GetBody() string { return c.Body }

// GetScore returns the current score of the Comment
func (c Comment) GetScore() int { return c.Score }

// IsSticky tells you if the comment has been stickied
func (c Comment) IsSticky() bool { return c.Stickied }

// IsRemoved tells you if the comment has been removed (only if you are mod)
func (c Comment) IsRemoved() bool { return c.Removed }

// IsApproved tells you if the comment has been approved (only if you are mod)
func (c Comment) IsApproved() bool { return c.Approved }

// IsAuthor tells you if the comment has been made by OP
func (c Comment) IsAuthor() bool { return c.IsSubmitter }

// GetURL returns the link to the Comment
func (c Comment) GetURL() string { return fmt.Sprintf("https://www.reddit.com%s", c.Permalink) }
