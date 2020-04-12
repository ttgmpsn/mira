package mira

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type transport struct {
	http.RoundTripper
	useragent string
}

// Any request headers can be modified here.
func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("User-Agent", t.useragent)
	return t.RoundTripper.RoundTrip(req)
}

// Taken mostly from https://github.com/jzelinskie/geddit/blob/master/oauth_session.go
// Thanks geddit :)

// NewOAuthSession creates a new session for those who want to log into a
// reddit account via OAuth.
func NewOAuthSession(creds Credentials) *Reddit {
	r := &Reddit{creds: creds}

	if len(r.creds.UserAgent) == 0 {
		r.creds.UserAgent = "unconfigured reddit bot using https://github.com/ttgmpsn/mira"
	}

	// Set OAuth config
	r.OAuthConfig = &oauth2.Config{
		ClientID:     r.creds.ClientID,
		ClientSecret: r.creds.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.reddit.com/api/v1/authorize",
			TokenURL: "https://www.reddit.com/api/v1/access_token",
		},
		RedirectURL: r.creds.RedirectURL,
	}
	// Inject our custom HTTP client so that a user-defined UA can
	// be passed during any authentication requests.
	c := &http.Client{}
	c.Transport = &transport{http.DefaultTransport, r.creds.UserAgent}
	r.ctx = context.WithValue(context.Background(), oauth2.HTTPClient, c)
	return r
}

// LoginAuth creates the required HTTP client with a new token.
// Creds are taken from the data provided to NewOAuthSession.
func (c *Reddit) LoginAuth() error {
	if len(c.creds.Username) == 0 || len(c.creds.Password) == 0 {
		return errors.New("no username or password provided to NewOAuthSession")
	}

	// Fetch OAuth token.
	t, err := c.OAuthConfig.PasswordCredentialsToken(c.ctx, c.creds.Username, c.creds.Password)
	if err != nil {
		return err
	}
	if !t.Valid() {
		msg := "Invalid OAuth token"
		if t != nil {
			if extra := t.Extra("error"); extra != nil {
				msg = fmt.Sprintf("%s: %s", msg, extra)
			}
		}
		return errors.New(msg)
	}
	c.Client = c.OAuthConfig.Client(c.ctx, t)
	return nil
}

// AuthCodeURL creates and returns an auth URL which contains an auth code.
func (c *Reddit) AuthCodeURL(state string, scopes []string) string {
	c.OAuthConfig.Scopes = scopes
	return c.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// CodeAuth creates and sets a token using an authentication code returned from AuthCodeURL.
// Note that Username & Password provided to NewOAuthSession are ignored.
func (c *Reddit) CodeAuth(code string) error {
	t, err := c.OAuthConfig.Exchange(c.ctx, code)
	if err != nil {
		return err
	}
	c.Client = c.OAuthConfig.Client(c.ctx, t)
	return nil
}

// SetToken manually assigns a token to the Reddit object.
// This is useful if you have your token information saved from a prior run.
func (c *Reddit) SetToken(t *oauth2.Token, scopes []string) {
	c.OAuthConfig.Scopes = scopes
	c.Client = c.OAuthConfig.Client(c.ctx, t)
}

// SetDefault gets sensible default values for streams.
func (c *Reddit) SetDefault() {
	c.Values = RedditVals{
		GetSubmissionFromCommentTries: 32,
	}
}
