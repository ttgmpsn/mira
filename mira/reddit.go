package mira

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

func (c *Reddit) Me() (Me, error) {
	target := RedditApiMe
	user := Me{}
	r, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return user, err
	}
	r.Header.Set("User-Agent", c.Creds.UserAgent)
	r.Header.Set("Authorization", "bearer "+c.Token)
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return user, err
	}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	json.Unmarshal(buf.Bytes(), &user)
	return user, nil
}

func (c *Reddit) GetUser(name string) (Redditor, error) {
	target := RedditOauth + "/user/" + name + "/about"
	user := Redditor{}
	r, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return user, err
	}
	r.Header.Set("User-Agent", c.Creds.UserAgent)
	r.Header.Set("Authorization", "bearer "+c.Token)
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return user, err
	}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	json.Unmarshal(buf.Bytes(), &user)
	return user, nil
}

func (c *Reddit) Submit(sr string, title string, text string) (Submission, error) {
	target := RedditOauth + "/api/submit"
	post := Submission{}
	form := url.Values{}
	form.Add("title", title)
	form.Add("sr", sr)
	form.Add("text", text)
	form.Add("kind", "self")
	form.Add("resubmit", "true")
	form.Add("api_type", "json")
	r, err := http.NewRequest("POST", target, strings.NewReader(form.Encode()))
	if err != nil {
		return post, err
	}
	r.Header.Set("User-Agent", c.Creds.UserAgent)
	r.Header.Set("Authorization", "bearer "+c.Token)
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return post, err
	}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	json.Unmarshal(buf.Bytes(), &post)
	return post, nil
}

func (c *Reddit) Reply(comment_id string, text string) (Comment, error) {
	target := RedditOauth + "/api/comment"
	comment := Comment{}
	form := url.Values{}
	form.Add("text", text)
	form.Add("thing_id", comment_id)
	form.Add("api_type", "json")
	r, err := http.NewRequest("POST", target, strings.NewReader(form.Encode()))
	if err != nil {
		return comment, err
	}
	r.Header.Set("User-Agent", c.Creds.UserAgent)
	r.Header.Set("Authorization", "bearer "+c.Token)
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return comment, err
	}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	json.Unmarshal(buf.Bytes(), &comment)
	return comment, nil
}

func (c *Reddit) Comment(submission_id, text string) (Comment, error) {
	target := RedditOauth + "/api/comment"
	comment := Comment{}
	form := url.Values{}
	form.Add("text", text)
	form.Add("thing_id", submission_id)
	form.Add("api_type", "json")
	r, err := http.NewRequest("POST", target, strings.NewReader(form.Encode()))
	if err != nil {
		return comment, err
	}
	r.Header.Set("User-Agent", c.Creds.UserAgent)
	r.Header.Set("Authorization", "bearer "+c.Token)
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return comment, err
	}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	json.Unmarshal(buf.Bytes(), &comment)
	return comment, nil
}

func (c *Reddit) DeleteComment(comment_id string) (error) {
	target := RedditOauth + "/api/del"
	form := url.Values{}
	form.Add("id", comment_id)
	form.Add("api_type", "json")
	r, err := http.NewRequest("POST", target, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	r.Header.Set("User-Agent", c.Creds.UserAgent)
	r.Header.Set("Authorization", "bearer "+c.Token)
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func (c *Reddit) EditComment(comment_id, text string) (Comment, error) {
	target := RedditOauth + "/api/editusertext"
	comment := Comment{}
	form := url.Values{}
	form.Add("text", text)
	form.Add("thing_id", comment_id)
	form.Add("api_type", "json")
	r, err := http.NewRequest("POST", target, strings.NewReader(form.Encode()))
	if err != nil {
		return comment, err
	}
	r.Header.Set("User-Agent", c.Creds.UserAgent)
	r.Header.Set("Authorization", "bearer "+c.Token)
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return comment, err
	}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	json.Unmarshal(buf.Bytes(), &comment)
	return comment, nil
}

func (c *Reddit) Compose(to, subject, text string) (error) {
	target := RedditOauth + "/api/compose"
	form := url.Values{}
	form.Add("subject", subject)
	form.Add("text", text)
	form.Add("to", to)
	form.Add("api_type", "json")	
	r, err := http.NewRequest("POST", target, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	r.Header.Set("User-Agent", c.Creds.UserAgent)
	r.Header.Set("Authorization", "bearer "+c.Token)
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func (c *Reddit) ReadMessage(message_id string) (error) {
	target := RedditOauth + "/api/read_message"
	form := url.Values{}
	form.Add("id", message_id)
	r, err := http.NewRequest("POST", target, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	r.Header.Set("User-Agent", c.Creds.UserAgent)
	r.Header.Set("Authorization", "bearer "+c.Token)
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func (c *Reddit) ReadAllMessages() (error) {
	target := RedditOauth + "/api/read_all_messages"
	r, err := http.NewRequest("POST", target, nil)
	if err != nil {
		return err
	}
	r.Header.Set("User-Agent", c.Creds.UserAgent)
	r.Header.Set("Authorization", "bearer "+c.Token)
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func (c *Reddit) ListUnreadMessages() (Listing, error) {
	target := RedditOauth + "/message/unread"
	list := Listing{}
	form := url.Values{}
	form.Add("mark", "true")
	r, err := http.NewRequest("GET", target, strings.NewReader(form.Encode()))
	if err != nil {
		return list, err
	}
	r.Header.Set("User-Agent", c.Creds.UserAgent)
	r.Header.Set("Authorization", "bearer "+c.Token)
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return list, err
	}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	json.Unmarshal(buf.Bytes(), &list)	
	return list, nil
}