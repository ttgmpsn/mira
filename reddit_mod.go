package mira

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ttgmpsn/mira/models"
)

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

// ModQueue returns the mod queue from the last queued object.
// Valid objects: Subreddit
func (c *Reddit) ModQueue(limit int) ([]models.Submission, error) {
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

	ret := []models.Submission{}
	for _, post := range list.Children {
		if p, ok := post.Data.(models.Submission); ok {
			ret = append(ret, p)
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

// Ban bans a redditor from last queued object.
// Valid objects: Subreddit
func (c *Reddit) Ban(redditor string, days int, context, message, reason string) error {
	subreddit, _, err := c.checkType(models.KSubreddit)
	if err != nil {
		return err
	}
	target := RedditOauth + "/r/" + subreddit + "/api/friend"
	_, err = c.MiraRequest("POST", target, map[string]string{
		"name":        redditor,
		"duration":    strconv.Itoa(days),
		"ban_context": context,
		"ban_message": message,
		"ban_reason":  reason,
		"note":        reason,
		"api_type":    "json",
		"type":        "banned",
	})
	return err
}

// GetModMailByID returns the ModMail Conversation for a given modmail ID
func (c *Reddit) GetModMailByID(conversationID string, markRead bool) (*models.NewModmailConversation, error) {
	target := RedditOauth + "/api/mod/conversations/" + conversationID
	ans, err := c.MiraRequest("GET", target, map[string]string{
		"markRead": strconv.FormatBool(markRead),
	})
	if err != nil {
		return nil, err
	}
	ret := &models.NewModmailConversation{}
	if err := json.Unmarshal([]byte(ans), ret); err != nil {
		return nil, err
	}
	return ret, nil
}
