package mira

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ttgmpsn/mira/models"
)

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

// Get submisssions from a subreddit up to a specified limit sorted by the given parameters
// and with specified anchor
//
// Sorting options: "hot", "new", "top", "rising", "controversial", "random"
//
// Limit is any numerical value, so 0 <= limit <= 100
//
// Anchor options are submissions full thing, for example: t3_bqqwm3
func (c *Reddit) getSubredditPostsAfter(sr string, last models.RedditID, limit int) ([]*models.Post, error) {
	target := RedditOauth + "/r/" + sr + "/new.json"
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

func (c *Reddit) getSubredditCommentsAfter(sr string, sort string, last models.RedditID, limit int) ([]*models.Comment, error) {
	target := RedditOauth + "/r/" + sr + "/comments.json"
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

// Wiki returns a wiki page from last queued object.
// Valid objects: Subreddit
func (c *Reddit) Wiki(page string) (*models.Wiki, error) {
	sr, ttype := c.getQueue()
	if ttype != models.KSubreddit {
		return nil, fmt.Errorf("'%s' type does not have an option for wiki", ttype)
	}

	target := RedditOauth + "/r/" + sr + "/wiki/" + page + ".json"
	ans, err := c.MiraRequest("GET", target, map[string]string{})
	if err != nil {
		return nil, err
	}
	ret := &models.Response{}
	if err := json.Unmarshal([]byte(ans), ret); err != nil {
		return nil, err
	}

	wiki, ok := ret.Data.(*models.Wiki)
	if !ok {
		return nil, fmt.Errorf("couldn't convert to Wiki struct")
	}

	return wiki, nil
}

// Stylesheet returns the stylesheet & images from last queued object.
// Valid objects: Subreddit
func (c *Reddit) Stylesheet() (*models.Stylesheet, error) {
	sr, ttype := c.getQueue()
	if ttype != models.KSubreddit {
		return nil, fmt.Errorf("'%s' type does not have an option for stylesheet", ttype)
	}

	target := RedditOauth + "/r/" + sr + "/about/stylesheet.json"
	ans, err := c.MiraRequest("GET", target, map[string]string{})
	if err != nil {
		return nil, err
	}
	ret := &models.Response{}
	if err := json.Unmarshal([]byte(ans), ret); err != nil {
		return nil, err
	}

	stylesheet, ok := ret.Data.(*models.Stylesheet)
	if !ok {
		return nil, fmt.Errorf("couldn't convert to Wiki struct")
	}

	return stylesheet, nil
}
