package mira

import (
	"container/ring"
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
// The fetch interval can be set via reddit.Config.CommentStreamInterval
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
// The fetch interval can be set via reddit.Config.PostStreamInterval
// Valid objects: Subreddit, (Redditor)
func (c *Reddit) StreamPosts() (*SubmissionStream, error) {
	name, ttype := c.getQueue()
	switch ttype {
	case models.KSubreddit:
		return c.streamSubredditPosts(name)
	/*case models.KRedditor:
	return c.streamRedditorPosts(name)*/
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
		sent := ring.New(100)
		for {
			select {
			case <-s.Close:
				return
			default:
			}
			comments, err := c.Subreddit(name).CommentsAfter("new", last, 100)
			if err != nil {
				close(sendC)
				return
			}
			for i := len(comments) - 1; i >= 0; i-- {
				if ringContains(sent, comments[i].GetID()) {
					continue
				}
				sendC <- comments[i]
				sent.Value = comments[i].GetID()
				sent = sent.Next()
			}
			if len(comments) == 0 {
				last = ""
			} else if len(comments) > 2 {
				last = comments[1].GetID()
			}
			time.Sleep(time.Duration(c.Config.CommentStreamInterval) * time.Second)
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
		sent := ring.New(100)
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
				if ringContains(sent, posts[i].GetID()) {
					continue
				}
				sendC <- posts[i]
				sent.Value = posts[i].GetID()
				sent = sent.Next()
			}
			if len(posts) == 0 {
				last = ""
			} else if len(posts) > 2 {
				last = posts[1].GetID()
			}
			time.Sleep(time.Duration(c.Config.PostStreamInterval) * time.Second)
		}
	}()
	return s, nil
}

func ringContains(r *ring.Ring, n models.RedditID) bool {
	ret := false
	r.Do(func(p interface{}) {
		if p == nil {
			return
		}
		i := p.(models.RedditID)
		if i == n {
			ret = true
		}
	})
	return ret
}
