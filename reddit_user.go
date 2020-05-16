package mira

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ttgmpsn/mira/models"
)

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

func (c *Reddit) getRedditorPostsAfter(user string, last models.RedditID, limit int) ([]*models.Post, error) {
	target := RedditOauth + "/u/" + user + "/submitted/new.json"
	list, err := c.miraRequestListing("GET", target, map[string]string{
		"limit":  strconv.Itoa(limit),
		"before": string(last),
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

func (c *Reddit) getRedditorCommentsAfter(user string, sort string, last models.RedditID, limit int) ([]*models.Comment, error) {
	target := RedditOauth + "/u/" + user + "/comments.json"
	list, err := c.miraRequestListing("GET", target, map[string]string{
		"sort":   sort,
		"limit":  strconv.Itoa(limit),
		"before": string(last),
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
	json.Unmarshal(ans, &ret)
	return ret, err
}
