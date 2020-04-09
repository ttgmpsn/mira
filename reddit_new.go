package mira

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ttgmpsn/mira/models"
)

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

// ModQueue returns the mod queue from the last queued object.
// Valid objects: Subreddit
func (c *Reddit) ModQueue(limit int) ([]*models.Submission, error) {
	sr, ttype := c.getQueue()
	if ttype != models.KSubreddit {
		return nil, fmt.Errorf("'%s' type does not have an option for modqueue", ttype)
	}

	target := RedditOauth + "/r/" + sr + "/about/modqueue.json"
	list, err := c.miraRequestListing("GET", target, map[string]string{
		"limit": strconv.Itoa(limit),
	})
	if err != nil {
		return nil, err
	}

	ret := []*models.Submission{}
	for _, post := range list.Children {
		if p, ok := post.Data.(models.Submission); ok {
			ret = append(ret, &p)
		}
	}

	return ret, nil
}

// ModLog returns the mod log from the last queued object.
// Valid objects: Subreddit
func (c *Reddit) ModLog(limit int, mod string) ([]*models.ModAction, error) {
	sr, ttype := c.getQueue()
	if ttype != models.KSubreddit {
		return nil, fmt.Errorf("'%s' type does not have an option for modlog", ttype)
	}

	target := RedditOauth + "/r/" + sr + "/about/log.json"
	list, err := c.miraRequestListing("GET", target, map[string]string{
		"limit": strconv.Itoa(limit),
		"mod":   mod,
	})
	if err != nil {
		return nil, err
	}

	ret := []*models.ModAction{}
	for _, post := range list.Children {
		if p, ok := post.Data.(*models.ModAction); ok {
			ret = append(ret, p)
		}
	}

	return ret, nil
}
