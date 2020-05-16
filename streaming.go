package mira

import (
	"fmt"
	"time"

	"github.com/ttgmpsn/mira/models"
)

// SubmissionStream has two objects: a channel "C" where you can receive
// Submissions (posts or comments), and a channel "close" - close that
// channel to stop receiving events. Please use the close channel appropriately
// or you'll be polling reddit non-stop!
type SubmissionStream struct {
	C     <-chan models.Submission
	Close chan struct{}
}

// StreamComments streams comments for the last queued object.
// Valid objects: Subreddit, (Redditor)
func (c *Reddit) StreamComments() (*SubmissionStream, error) {
	name, ttype := c.getQueue()
	switch ttype {
	case models.KSubreddit:
		return c.streamSubredditComments(name)
	/*case models.KRedditor:
	return c.streamRedditorComments(name)*/
	default:
		return nil, fmt.Errorf("'%s' type does not have an option to stream comments", ttype)
	}
}

// StreamPosts streams posts for the last queued object.
// Valid objects: Subreddit, (Redditor)
func (c *Reddit) StreamPosts() (*SubmissionStream, error) {
	name, ttype := c.getQueue()
	switch ttype {
	case models.KSubreddit:
		return c.streamSubredditPosts(name)
	/*case models.KRedditor:
	return c.streamRedditorComments(name)*/
	default:
		return nil, fmt.Errorf("'%s' type does not have an option to stream posts", ttype)
	}
}

func (c *Reddit) streamSubredditComments(name string) (*SubmissionStream, error) {
	sendC := make(chan models.Submission, 100)
	s := &SubmissionStream{
		C:     sendC,
		Close: make(chan struct{}),
	}
	_, err := c.Subreddit(name).Posts("new", "all", 1)
	if err != nil {
		return nil, err
	}
	var last models.RedditID
	go func() {
		for {
			select {
			case <-s.Close:
				close(sendC)
				return
			default:
			}
			comments, err := c.Subreddit(name).CommentsAfter("new", last, 100)
			if err != nil {
				close(sendC)
				return
			}
			for i := len(comments) - 1; i >= 0; i-- {
				sendC <- comments[i]
				last = comments[i].GetID()
			}
			time.Sleep(45 * time.Second)
		}
	}()
	return s, nil
}

func (c *Reddit) streamSubredditPosts(name string) (*SubmissionStream, error) {
	sendC := make(chan models.Submission, 100)
	s := &SubmissionStream{
		C:     sendC,
		Close: make(chan struct{}),
	}
	_, err := c.Subreddit(name).Posts("new", "all", 1)
	if err != nil {
		return nil, err
	}
	var last models.RedditID
	go func() {
		for {
			select {
			case <-s.Close:
				return
			default:
			}
			posts, err := c.Subreddit(name).PostsAfter(last, 100)
			if err != nil {
				close(sendC)
				return
			}
			for i := len(posts) - 1; i >= 0; i-- {
				sendC <- posts[i]
				last = posts[i].GetID()
			}
			time.Sleep(45 * time.Second)
		}
	}()
	return s, nil
}
