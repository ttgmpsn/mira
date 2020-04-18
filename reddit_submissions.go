package mira

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/ttgmpsn/mira/models"
)

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
