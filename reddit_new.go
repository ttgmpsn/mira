package mira

import (
	"encoding/json"
	"fmt"

	"github.com/ttgmpsn/mira/models"
)

// Wiki returns a wiki page from last queued object.
// Valid objects: Subreddit
func (c *Reddit) Wiki(page string) (*models.WikiData, error) {
	sr, ttype := c.getQueue()
	if ttype != "subreddit" {
		return nil, fmt.Errorf("'%s' type does not have an option for mod log", ttype)
	}

	target := RedditOauth + "/r/" + sr + "/wiki/" + page + ".json"
	ans, err := c.MiraRequest("GET", target, map[string]string{})
	ret := &models.Wiki{}
	json.Unmarshal(ans, ret)
	return &ret.Data, err
}
