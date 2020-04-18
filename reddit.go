package mira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/ttgmpsn/mira/models"
)

var (
	requestMutex *sync.RWMutex = &sync.RWMutex{}
	queueMutex   *sync.RWMutex = &sync.RWMutex{}
)

// MiraRequest can be used to make custom requests to the reddit API.
func (c *Reddit) MiraRequest(method string, target string, payload map[string]string) ([]byte, error) {
	values := "?"
	for i, v := range payload {
		v = url.QueryEscape(v)
		values += fmt.Sprintf("%s=%s&", i, v)
	}
	values = values[:len(values)-1]
	r, err := http.NewRequest(method, target+values, nil)
	if err != nil {
		return nil, err
	}
	requestMutex.Lock()
	response, err := c.Client.Do(r)
	requestMutex.Unlock()
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

// Me Redditor queues up the next action to be about the logged in user.
func (c *Reddit) Me() *Reddit {
	c.addQueue("", "me")
	return c
}

// Subreddit Redditor queues up the next action to be about one or multuple Subreddits.
func (c *Reddit) Subreddit(name ...string) *Reddit {
	c.addQueue(strings.Join(name, "+"), models.KSubreddit)
	return c
}

// Post queues up the next action to be about a certain Post.
func (c *Reddit) Post(name string) *Reddit {
	c.addQueue(name, models.KPost)
	return c
}

// Comment queues up the next action to be about a certain comment.
func (c *Reddit) Comment(name string) *Reddit {
	c.addQueue(name, models.KComment)
	return c
}

// Redditor queues up the next action to be about a certain Redditor.
func (c *Reddit) Redditor(name string) *Reddit {
	c.addQueue(name, models.KRedditor)
	return c
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
func (c *Reddit) PostsAfter(last string, limit int) ([]*models.Post, error) {
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
func (c *Reddit) CommentsAfter(sort string, last string, limit int) ([]*models.Comment, error) {
	name, ttype := c.getQueue()
	switch ttype {
	case "subreddit":
		return c.getSubredditCommentsAfter(name, sort, last, limit)
	case "redditor":
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

func (c *Reddit) addQueue(name string, ttype models.RedditKind) {
	queueMutex.Lock()
	defer queueMutex.Unlock()
	c.Chain = append(c.Chain, chainVals{Name: name, Type: ttype})
}

func (c *Reddit) getQueue() (string, models.RedditKind) {
	queueMutex.Lock()
	defer queueMutex.Unlock()
	if len(c.Chain) < 1 {
		return "", ""
	}
	defer func() { c.Chain = c.Chain[1:] }()
	return c.Chain[0].Name, c.Chain[0].Type
}

func findElem(elem models.RedditKind, arr []models.RedditKind) bool {
	for _, v := range arr {
		if elem == v {
			return true
		}
	}
	return false
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
