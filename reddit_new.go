package mira

import (
	"encoding/json"
	"fmt"

	"github.com/ttgmpsn/mira/models"
)

// Wiki returns a wiki page from last queued object.
// Valid objects: Subreddit
func (c *Reddit) Wiki(page string) (*models.Wiki, error) {
	sr, ttype := c.getQueue()
	if ttype != "subreddit" {
		return nil, fmt.Errorf("'%s' type does not have an option for mod log", ttype)
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
