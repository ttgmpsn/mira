package mira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ttgmpsn/mira/models"
)

// MiraRequest can be used to make custom requests to the reddit API.
func (c *Reddit) MiraRequest(method string, target string, payload map[string]string) ([]byte, error) {
	values := url.Values{}
	for i, v := range payload {
		values.Set(i, v)
	}

	var r *http.Request
	var err error
	if method == "GET" {
		values := fmt.Sprintf("?%s", values.Encode())
		r, err = http.NewRequest(method, target+values, nil)
	} else {
		r, err = http.NewRequest(method, target, strings.NewReader(values.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("Content-Length", strconv.Itoa(len(values.Encode())))
	}
	if err != nil {
		return nil, err
	}
	response, err := c.Client.Do(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	data := buf.Bytes()
	if err := findRedditError(data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Reddit) miraRequestListing(method string, target string, payload map[string]string) (*models.Listing, error) {
	ans, err := c.MiraRequest(method, target, payload)
	if err != nil {
		return nil, err
	}
	ret := &models.Response{}
	if err := json.Unmarshal([]byte(ans), ret); err != nil {
		return nil, err
	}
	list, ok := ret.Data.(*models.Listing)
	if !ok {
		return nil, fmt.Errorf("couldn't convert to Listing struct. Data has Kind '%s'", ret.Kind)
	}
	return list, nil
}

// Me queues up the next action to be about the logged in user.
func (c *Reddit) Me() *Reddit {
	return c.addQueue("", "me")
}

// Subreddit queues up the next action to be about one or multuple Subreddits.
func (c *Reddit) Subreddit(name ...string) *Reddit {
	return c.addQueue(strings.Join(name, "+"), models.KSubreddit)
}

// Post queues up the next action to be about a certain Post.
func (c *Reddit) Post(name string) *Reddit {
	return c.addQueue(name, models.KPost)
}

// Comment queues up the next action to be about a certain comment.
func (c *Reddit) Comment(name string) *Reddit {
	return c.addQueue(name, models.KComment)
}

// Redditor queues up the next action to be about a certain Redditor.
func (c *Reddit) Redditor(name string) *Reddit {
	return c.addQueue(name, models.KRedditor)
}

// Posts gets posts for the last queued object.
// Valid objects: Subreddit, Redditor
func (c *Reddit) Posts(sort string, tdur string, limit int) ([]*models.Post, error) {
	name, ttype := c.getQueue()
	switch ttype {
	case models.KSubreddit:
		return c.getSubredditPosts(name, sort, tdur, limit)
	case models.KRedditor:
		return c.getRedditorPosts(name, sort, tdur, limit)
	default:
		return nil, fmt.Errorf("'%s' type does not have an option for posts", ttype)
	}
}

// PostsAfter gets posts for the last queued object after a given item.
// Valid objects: Subreddit, Redditor
func (c *Reddit) PostsAfter(last models.RedditID, limit int) ([]*models.Post, error) {
	name, ttype := c.getQueue()
	switch ttype {
	case models.KSubreddit:
		return c.getSubredditPostsAfter(name, last, limit)
	case models.KRedditor:
		return c.getRedditorPostsAfter(name, last, limit)
	default:
		return nil, fmt.Errorf("'%s' type does not have an option for postsafter", ttype)
	}
}

// Comments gets comments for the last queued object.
// Valid objects: Subreddit, Post, Redditor
func (c *Reddit) Comments(sort string, tdur string, limit int) ([]*models.Comment, error) {
	name, ttype := c.getQueue()
	switch ttype {
	case models.KSubreddit:
		return c.getSubredditComments(name, sort, tdur, limit)
	case models.KPost:
		comments, err := c.getPostComments(models.RedditID(name), sort, tdur, limit)
		if err != nil {
			return nil, err
		}
		return comments, nil
	case models.KRedditor:
		return c.getRedditorComments(name, sort, tdur, limit)
	default:
		return nil, fmt.Errorf("'%s' type does not have an option for comments", ttype)
	}
}

// Info returns general information about the queued object as a mira.Interface.
func (c *Reddit) Info() (models.RedditThing, error) {
	name, ttype := c.getQueue()
	switch ttype {
	case "me":
		return c.getMe()
	case models.KPost:
		return c.getPost(models.RedditID(name))
	case models.KComment:
		return c.getComment(models.RedditID(name))
	case models.KSubreddit:
		return c.getSubreddit(name)
	case models.KRedditor:
		return c.getUser(name)
	default:
		return nil, fmt.Errorf("returning type is not defined")
	}
}

// CommentsAfter gets comments for the last queued object after a given item.
// Valid objects: Subreddit, Redditor
func (c *Reddit) CommentsAfter(sort string, last models.RedditID, limit int) ([]*models.Comment, error) {
	name, ttype := c.getQueue()
	switch ttype {
	case models.KSubreddit:
		return c.getSubredditCommentsAfter(name, sort, last, limit)
	case models.KRedditor:
		return c.getRedditorCommentsAfter(name, sort, last, limit)
	default:
		return nil, fmt.Errorf("'%s' type does not have an option for comments", ttype)
	}
}

func (c *Reddit) checkType(rtype ...models.RedditKind) (string, models.RedditKind, error) {
	name, ttype := c.getQueue()
	if name == "" {
		return "", "", fmt.Errorf("identifier is empty")
	}
	if !findElem(ttype, rtype) {
		return "", "", fmt.Errorf("the passed type is not a valid type for this call | expected: %s", rtype)
	}
	return name, ttype, nil
}

func (c *Reddit) addQueue(name string, ttype models.RedditKind) *Reddit {
	c.chain <- &chainVals{Name: name, Type: ttype}
	return c
}

func (c *Reddit) getQueue() (string, models.RedditKind) {
	next := <-c.chain
	return next.Name, next.Type
}

func findElem(elem models.RedditKind, arr []models.RedditKind) bool {
	for _, v := range arr {
		if elem == v {
			return true
		}
	}
	return false
}

// Submissions gets submissions for the last queued object.
// Valid objects: Redditor
func (c *Reddit) Submissions(limit int) ([]models.Submission, error) {
	name, ttype := c.getQueue()
	switch ttype {
	case models.KRedditor:
		return c.getRedditorSubmissions(name, limit)
	default:
		return nil, fmt.Errorf("'%s' type does not have an option for submissions", ttype)
	}
}

// SubmissionsAfter gets submissions for the last queued object after a given item.
// Valid objects: Redditor
func (c *Reddit) SubmissionsAfter(last models.RedditID, limit int) ([]models.Submission, error) {
	name, ttype := c.getQueue()
	switch ttype {
	case models.KRedditor:
		return c.getRedditorSubmissionsAfter(name, last, limit)
	default:
		return nil, fmt.Errorf("'%s' type does not have an option for submissionsafter", ttype)
	}
}

// RedditErr is an error returned from the Reddit API.
type RedditErr struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func findRedditError(data []byte) error {
	object := &RedditErr{}
	json.Unmarshal(data, object)
	if object.Message != "" || object.Error != "" {
		return fmt.Errorf("%s | error code: %s", object.Message, object.Error)
	}
	return nil
}
