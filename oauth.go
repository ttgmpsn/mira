package mira

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/oauth2"
)

type transport struct {
	http.RoundTripper
	useragent string
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("User-Agent", t.useragent)
	return t.RoundTripper.RoundTrip(req)
}

// Taken mostly from https://github.com/jzelinskie/geddit/blob/master/oauth_session.go
// Thanks geddit :)

// newOAuthSession creates a new session for those who want to log into a
// reddit account via OAuth.
func newOAuthSession(creds Credentials) *Reddit {
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
// Creds are taken from the data provided to Init.
// Tokens are refreshed automatically shortly before the session runs out.
func (c *Reddit) LoginAuth() error {
	if len(c.creds.Username) == 0 || len(c.creds.Password) == 0 {
		return errors.New("no username or password provided to Init")
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

	ts := &loginAuthRefreshTokenSource{
		t: t,
		c: c,
	}

	c.Client = oauth2.NewClient(c.ctx, ts)
	return nil
}

// AuthCodeURL creates and returns an auth URL which contains an auth code.
func (c *Reddit) AuthCodeURL(state string, scopes []string) string {
	c.OAuthConfig.Scopes = scopes
	return c.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline, oauth2.SetAuthURLParam("duration", "permanent"))
}

// CodeAuth creates and sets a token using an authentication code returned from AuthCodeURL.
// Note that Username & Password provided to Init are ignored.
// You can optionally pass a TokenNotifyFunc to get notified when the token changes (i.e. to
// store it into a database). Pass nil if you do not want to use this.
func (c *Reddit) CodeAuth(code string, f TokenNotifyFunc) error {
	t, err := c.OAuthConfig.Exchange(c.ctx, code)
	if err != nil {
		return err
	}

	return c.SetToken(t, c.OAuthConfig.Scopes, f)
}

// SetToken manually assigns a token to the Reddit object.
// This is useful if you have your token information saved from a prior run.
// You can optionally pass a TokenNotifyFunc to get notified when the token changes (i.e. to
// store it into a database). Pass nil if you do not want to use this.
func (c *Reddit) SetToken(t *oauth2.Token, scopes []string, f TokenNotifyFunc) error {
	c.OAuthConfig.Scopes = scopes

	if f == nil {
		f = func(t *oauth2.Token) error { return nil }
	}

	err := f(t)
	if err != nil {
		return err
	}

	nrts := &NotifyRefreshTokenSource{
		new: c.OAuthConfig.TokenSource(c.ctx, t),
		t:   t,
		f:   f,
	}

	c.Client = oauth2.NewClient(c.ctx, nrts)
	return nil
}

// TokenNotifyFunc is a function that accepts an oauth2 Token upon refresh, and
// returns an error if it should not be used. Use this to cache Refresh Token
// if you want to (you'll most likely want to).
// Taken from https://github.com/golang/oauth2/issues/84#issuecomment-332517319
type TokenNotifyFunc func(*oauth2.Token) error

// NotifyRefreshTokenSource is essentially oauth2.ResuseTokenSource with TokenNotifyFunc added.
type NotifyRefreshTokenSource struct {
	new oauth2.TokenSource
	mu  sync.Mutex // guards t
	t   *oauth2.Token
	f   TokenNotifyFunc // called when token refreshed so new refresh token can be persisted
}

// Token returns the current token if it's still valid, else will
// refresh the current token (using r.Context for HTTP client
// information) and return the new one.
func (s *NotifyRefreshTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.t.Valid() {
		return s.t, nil
	}
	t, err := s.new.Token()
	if err != nil {
		return nil, err
	}
	s.t = t
	return t, s.f(t)
}

// loginAuthRefreshTokenSource is essentially `oauth2.ResuseTokenSource`
// that re-logins every time the token runs out.
type loginAuthRefreshTokenSource struct {
	mu sync.Mutex // guards t
	t  *oauth2.Token
	c  *Reddit
}

// Token returns the current token if it's still valid, else will
// refresh the current token (using r.Context for HTTP client
// information) and return the new one.
func (s *loginAuthRefreshTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.t.Valid() {
		return s.t, nil
	}

	// Fetch OAuth token.
	t, err := s.c.OAuthConfig.PasswordCredentialsToken(s.c.ctx, s.c.creds.Username, s.c.creds.Password)
	if err != nil {
		return nil, err
	}
	if !t.Valid() {
		msg := "Invalid OAuth token"
		if t != nil {
			if extra := t.Extra("error"); extra != nil {
				msg = fmt.Sprintf("%s: %s", msg, extra)
			}
		}
		return nil, errors.New(msg)
	}
	s.t = t
	return t, nil
}
