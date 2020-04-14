package mira

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

// SubmissionInfo returns general information about the queued submission.
func (c *Reddit) SubmissionInfo() (models.Submission, error) {
	name, ttype := c.getQueue()
	switch ttype {
	case models.KPost:
		return c.getPost(models.RedditID(name))
	case models.KComment:
		return c.getComment(models.RedditID(name))
	default:
		return nil, fmt.Errorf("returning type is not defined")
	}
}

// SubmissionInfoID returns general information about the submission ID.
func (c *Reddit) SubmissionInfoID(name models.RedditID) (models.Submission, error) {
	switch name.Type() {
	case models.KPost:
		return c.getPost(name)
	case models.KComment:
		return c.getComment(name)
	default:
		return nil, fmt.Errorf("returning type is not defined")
	}
}

func (c *Reddit) getMe() (*models.Me, error) {
	target := RedditOauth + "/api/v1/me"
	ans, err := c.MiraRequest("GET", target, nil)
	if err != nil {
		return nil, err
	}
	ret := &models.Me{}
	if err := json.Unmarshal([]byte(ans), ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *Reddit) getPost(id models.RedditID) (*models.Post, error) {
	target := RedditOauth + "/api/info.json"
	list, err := c.miraRequestListing("GET", target, map[string]string{
		"id": string(id),
	})
	if err != nil {
		return nil, err
	}

	if len(list.Children) < 1 {
		return nil, fmt.Errorf("no results")
	}

	post, ok := list.Children[0].Data.(*models.Post)
	if !ok {
		return nil, fmt.Errorf("provided ID '%s' is no valid post", id)
	}

	return post, nil
}

func (c *Reddit) getComment(id models.RedditID) (*models.Comment, error) {
	target := RedditOauth + "/api/info.json"
	list, err := c.miraRequestListing("GET", target, map[string]string{
		"id": string(id),
	})
	if err != nil {
		return nil, err
	}

	if len(list.Children) < 1 {
		return nil, fmt.Errorf("no results")
	}

	comment, ok := list.Children[0].Data.(*models.Comment)
	if !ok {
		return nil, fmt.Errorf("provided ID '%s' is no valid comment", id)
	}

	return comment, nil
}

// GetParentPost returns the Post ID for the last queued object.
// Valid objects: Comment
func (c *Reddit) GetParentPost() (models.RedditID, error) {
	name, _, err := c.checkType(models.KComment)
	if err != nil {
		return "", err
	}
	info, err := c.getComment(models.RedditID(name))
	if err != nil {
		return "", err
	}
	return info.LinkID, nil
}

func (c *Reddit) getUser(name string) (*models.Redditor, error) {
	target := RedditOauth + "/user/" + name + "/about"
	ans, err := c.MiraRequest("GET", target, nil)
	if err != nil {
		return nil, err
	}
	ret := &models.RedditElement{}
	if err := json.Unmarshal([]byte(ans), ret); err != nil {
		return nil, err
	}
	user, ok := ret.Data.(*models.Redditor)
	if !ok {
		return nil, fmt.Errorf("couldn't convert to Redditor struct. Data has Kind '%s'", ret.Kind)
	}
	return user, nil
}

func (c *Reddit) getSubreddit(name string) (*models.Subreddit, error) {
	target := RedditOauth + "/r/" + name + "/about"
	ans, err := c.MiraRequest("GET", target, nil)
	if err != nil {
		return nil, err
	}
	ret := &models.RedditElement{}
	if err := json.Unmarshal([]byte(ans), ret); err != nil {
		return nil, err
	}
	sub, ok := ret.Data.(*models.Subreddit)
	if !ok {
		return nil, fmt.Errorf("couldn't convert to Subreddit struct. Data has Kind '%s'", ret.Kind)
	}
	return sub, nil
}

func (c *Reddit) getRedditorPosts(user string, sort string, tdur string, limit int) ([]*models.Post, error) {
	target := RedditOauth + "/u/" + user + "/submitted/" + sort + ".json"
	list, err := c.miraRequestListing("GET", target, map[string]string{
		"limit": strconv.Itoa(limit),
		"t":     tdur,
	})
	if err != nil {
		return nil, err
	}

	ret := []*models.Post{}
	for _, post := range list.Children {
		if p, ok := post.Data.(*models.Post); ok {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

func (c *Reddit) getRedditorPostsAfter(user string, last string, limit int) ([]*models.Post, error) {
	target := RedditOauth + "/u/" + user + "/submitted/new.json"
	list, err := c.miraRequestListing("GET", target, map[string]string{
		"limit":  strconv.Itoa(limit),
		"before": last,
	})
	if err != nil {
		return nil, err
	}

	ret := []*models.Post{}
	for _, post := range list.Children {
		if p, ok := post.Data.(*models.Post); ok {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

// Get submisssions from a subreddit up to a specified limit sorted by the given parameter
//
// Sorting options: "hot", "new", "top", "rising", "controversial", "random"
//
// Time options: "all", "year", "month", "week", "day", "hour"
//
// Limit is any numerical value, so 0 <= limit <= 100
func (c *Reddit) getSubredditPosts(sr string, sort string, tdur string, limit int) ([]*models.Post, error) {
	target := RedditOauth + "/r/" + sr + "/" + sort + ".json"
	list, err := c.miraRequestListing("GET", target, map[string]string{
		"limit": strconv.Itoa(limit),
		"t":     tdur,
	})
	if err != nil {
		return nil, err
	}

	ret := []*models.Post{}
	for _, post := range list.Children {
		if p, ok := post.Data.(*models.Post); ok {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

func (c *Reddit) getSubredditComments(sr string, sort string, tdur string, limit int) ([]*models.Comment, error) {
	target := RedditOauth + "/r/" + sr + "/comments.json"
	list, err := c.miraRequestListing("GET", target, map[string]string{
		"sort":  sort,
		"limit": strconv.Itoa(limit),
		"t":     tdur,
	})
	if err != nil {
		return nil, err
	}

	ret := []*models.Comment{}
	for _, comment := range list.Children {
		if c, ok := comment.Data.(*models.Comment); ok {
			ret = append(ret, c)
		}
	}

	return ret, nil
}

func (c *Reddit) getRedditorComments(user string, sort string, tdur string, limit int) ([]*models.Comment, error) {
	target := RedditOauth + "/u/" + user + "/comments.json"
	list, err := c.miraRequestListing("GET", target, map[string]string{
		"sort":  sort,
		"limit": strconv.Itoa(limit),
		"t":     tdur,
	})
	if err != nil {
		return nil, err
	}

	ret := []*models.Comment{}
	for _, comment := range list.Children {
		if c, ok := comment.Data.(*models.Comment); ok {
			ret = append(ret, c)
		}
	}

	return ret, nil
}

func (c *Reddit) getRedditorCommentsAfter(user string, sort string, last string, limit int) ([]*models.Comment, error) {
	target := RedditOauth + "/u/" + user + "/comments.json"
	list, err := c.miraRequestListing("GET", target, map[string]string{
		"sort":   sort,
		"limit":  strconv.Itoa(limit),
		"before": last,
	})
	if err != nil {
		return nil, err
	}

	ret := []*models.Comment{}
	for _, comment := range list.Children {
		if c, ok := comment.Data.(*models.Comment); ok {
			ret = append(ret, c)
		}
	}

	return ret, nil
}

func (c *Reddit) getPostComments(postID models.RedditID, sort string, tdur string, limit int) ([]*models.Comment, error) {
	if postID.Type() != models.KPost {
		return nil, errors.New("the passed ID is not a post")
	}
	target := fmt.Sprintf("%s/comments/%s", RedditOauth, postID[3:])
	ans, err := c.MiraRequest("GET", target, map[string]string{
		"sort":     sort,
		"limit":    strconv.Itoa(limit),
		"showmore": strconv.FormatBool(true),
		"t":        tdur,
	})
	if err != nil {
		return nil, err
	}

	rets := []*models.Response{}
	if err = json.Unmarshal(ans, &rets); err != nil {
		return nil, err
	}
	if len(rets) < 2 {
		// not two elements --> no comments
		return []*models.Comment{}, nil
	}
	list, ok := rets[1].Data.(*models.Listing)
	if !ok {
		return nil, fmt.Errorf("couldn't convert to Listing struct. Data has Kind '%s'", rets[1].Kind)
	}
	ret := []*models.Comment{}
	for _, comment := range list.Children {
		if c, ok := comment.Data.(*models.Comment); ok {
			ret = append(ret, c)
		}
	}
	return ret, nil
}

// Get submisssions from a subreddit up to a specified limit sorted by the given parameters
// and with specified anchor
//
// Sorting options: "hot", "new", "top", "rising", "controversial", "random"
//
// Limit is any numerical value, so 0 <= limit <= 100
//
// Anchor options are submissions full thing, for example: t3_bqqwm3
func (c *Reddit) getSubredditPostsAfter(sr string, last string, limit int) ([]*models.Post, error) {
	target := RedditOauth + "/r/" + sr + "/new.json"
	list, err := c.miraRequestListing("GET", target, map[string]string{
		"limit":  strconv.Itoa(limit),
		"before": last,
	})
	if err != nil {
		return nil, err
	}

	ret := []*models.Post{}
	for _, post := range list.Children {
		if p, ok := post.Data.(*models.Post); ok {
			ret = append(ret, p)
		}
	}

	return ret, nil
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

func (c *Reddit) getSubredditCommentsAfter(sr string, sort string, last string, limit int) ([]*models.Comment, error) {
	target := RedditOauth + "/r/" + sr + "/comments.json"
	list, err := c.miraRequestListing("GET", target, map[string]string{
		"sort":   sort,
		"limit":  strconv.Itoa(limit),
		"before": last,
	})
	if err != nil {
		return nil, err
	}

	ret := []*models.Comment{}
	for _, comment := range list.Children {
		if c, ok := comment.Data.(*models.Comment); ok {
			ret = append(ret, c)
		}
	}

	return ret, nil
}

// Submit submits a new Post to the last queued object.
// Valid objects: Subreddit
func (c *Reddit) Submit(title string, text string) (*models.Post, error) {
	ret := &models.Post{}
	name, _, err := c.checkType(models.KSubreddit)
	if err != nil {
		return nil, err
	}
	target := RedditOauth + "/api/submit"
	ans, err := c.MiraRequest("POST", target, map[string]string{
		"title":    title,
		"sr":       name,
		"text":     text,
		"kind":     "self",
		"resubmit": "true",
		"api_type": "json",
	})
	// :TODO: check reply type
	fmt.Println(string(ans))
	json.Unmarshal(ans, ret)
	return ret, err
}

// Reply adds a comment to the last queued object.
// Valid objects: Comment, Post
func (c *Reddit) Reply(text string) (*models.ActionResponse, error) {
	name, _, err := c.checkType(models.KComment, models.KPost)
	if err != nil {
		return nil, err
	}
	return c.ReplyWithID(name, text)
}

// ReplyWithID adds a comment to the given thing id, without it needing to be queued up.
func (c *Reddit) ReplyWithID(name, text string) (*models.ActionResponse, error) {
	ret := &models.ActionResponse{}
	target := RedditOauth + "/api/comment"
	ans, err := c.MiraRequest("POST", target, map[string]string{
		"text":     text,
		"thing_id": name,
		"api_type": "json",
	})
	json.Unmarshal(ans, ret)
	return ret, err
}

// Delete the last queued object.
// Valid objects: Comment, Post
func (c *Reddit) Delete() error {
	name, _, err := c.checkType(models.KComment, models.KPost)
	if err != nil {
		return err
	}
	target := RedditOauth + "/api/del"
	_, err = c.MiraRequest("POST", target, map[string]string{
		"id":       name,
		"api_type": "json",
	})
	return err
}

// Approve the last queued object.
// Valid objects: Comment, Post
func (c *Reddit) Approve() error {
	name, _, err := c.checkType(models.KComment, models.KPost)
	if err != nil {
		return err
	}
	target := RedditOauth + "/api/approve"
	_, err = c.MiraRequest("POST", target, map[string]string{
		"id":       name,
		"api_type": "json",
	})
	return err
}

// Remove mod-removes the last queued object. To remove own comments,
// please use Delete()
// Valid objects: Comment, Post
func (c *Reddit) Remove(spam bool) error {
	name, _, err := c.checkType(models.KComment, models.KPost)
	if err != nil {
		return err
	}
	target := RedditOauth + "/api/remove"
	_, err = c.MiraRequest("POST", target, map[string]string{
		"id":       name,
		"spam":     strconv.FormatBool(spam),
		"api_type": "json",
	})
	return err
}

// Distinguish the last queued object.
// Valid objects: Comment
func (c *Reddit) Distinguish(how string, sticky bool) error {
	name, _, err := c.checkType(models.KComment)
	if err != nil {
		return err
	}
	target := RedditOauth + "/api/distinguish"
	_, err = c.MiraRequest("POST", target, map[string]string{
		"id":       name,
		"how":      how,
		"sticky":   strconv.FormatBool(sticky),
		"api_type": "json",
	})
	return err
}

// Edit the last queued object.
// Valid objects: Comment, Post
func (c *Reddit) Edit(text string) (*models.Comment, error) {
	ret := &models.Comment{}
	name, _, err := c.checkType(models.KComment, models.KPost)
	if err != nil {
		return nil, err
	}
	target := RedditOauth + "/api/editusertext"
	ans, err := c.MiraRequest("POST", target, map[string]string{
		"text":     text,
		"thing_id": name,
		"api_type": "json",
	})
	// :TODO: check reply type
	fmt.Println(string(ans))
	json.Unmarshal(ans, ret)
	return ret, err
}

// Compose writes a private message to the last queued object.
// Valid objects: Redditor
func (c *Reddit) Compose(subject, text string) error {
	name, _, err := c.checkType(models.KRedditor)
	if err != nil {
		return err
	}
	target := RedditOauth + "/api/compose"
	_, err = c.MiraRequest("POST", target, map[string]string{
		"subject":  subject,
		"text":     text,
		"to":       name,
		"api_type": "json",
	})
	return err
}

// ReadMessage marks a message for the last queued object as read.
// Valid objects: Me
func (c *Reddit) ReadMessage(messageID string) error {
	_, _, err := c.checkType("me")
	if err != nil {
		return err
	}
	target := RedditOauth + "/api/read_message"
	_, err = c.MiraRequest("POST", target, map[string]string{
		"id": messageID,
	})
	return err
}

// ReadAllMessages marks all message for the last queued object as read.
// Valid objects: Me
func (c *Reddit) ReadAllMessages() error {
	_, _, err := c.checkType("me")
	if err != nil {
		return err
	}
	target := RedditOauth + "/api/read_all_messages"
	_, err = c.MiraRequest("POST", target, nil)
	return err
}

// ListUnreadMessages for the last queued object.
// Valid objects: Me
func (c *Reddit) ListUnreadMessages() ([]*models.Comment, error) {
	_, _, err := c.checkType("me")
	if err != nil {
		return nil, err
	}
	target := RedditOauth + "/message/unread"
	ans, err := c.MiraRequest("GET", target, map[string]string{
		"mark": "false",
	})
	ret := []*models.Comment{}
	// :TODO: check reply type
	fmt.Println(string(ans))
	json.Unmarshal(ans, ret)
	return ret, err
}

// UpdateSidebar of the last queued object.
// Valid objects: Subreddit
func (c *Reddit) UpdateSidebar(text string) error {
	name, _, err := c.checkType(models.KSubreddit)
	if err != nil {
		return err
	}
	target := RedditOauth + "/api/site_admin"
	_, err = c.MiraRequest("POST", target, map[string]string{
		"sr":          name,
		"name":        "None",
		"description": text,
		"title":       name,
		"wikimode":    "anyone",
		"link_type":   "any",
		"type":        "public",
		"api_type":    "json",
	})
	return err
}

// UserFlair assigns a specific flair to a user on the last queued object.
// Valid objects: Subreddit
func (c *Reddit) UserFlair(user, text string) error {
	name, _, err := c.checkType(models.KSubreddit)
	if err != nil {
		return err
	}
	target := RedditOauth + "/r/" + name + "/api/flair"
	_, err = c.MiraRequest("POST", target, map[string]string{
		"name":     user,
		"text":     text,
		"api_type": "json",
	})
	return err
}

// SelectFlair for the last queued object.
// Valid objects: Post
func (c *Reddit) SelectFlair(text string) error {
	name, _, err := c.checkType(models.KPost)
	if err != nil {
		return err
	}
	target := RedditOauth + "/api/selectflair"
	_, err = c.MiraRequest("POST", target, map[string]string{
		"link":     name,
		"text":     text,
		"api_type": "json",
	})
	return err
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
