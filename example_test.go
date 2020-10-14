package mira_test

import (
	"fmt"
	"time"

	"github.com/ttgmpsn/mira"
	miramodels "github.com/ttgmpsn/mira/models"
	"golang.org/x/oauth2"
)

func ExampleReddit_LoginAuth() {
	reddit := mira.Init(mira.Credentials{
		ClientID:     "clientid",
		ClientSecret: "clientsecret",
		Username:     "reddit_username",
		Password:     "topsecretpassword",
		UserAgent:    "MIRA LoginAuth Example v0",
	})

	err := reddit.LoginAuth()
	if err != nil {
		panic(err)
	}

	// You can now use reddit as an authenticated user:
	rMeObj, err := reddit.Me().Info()
	if err != nil {
		panic(err)
	}
	rMe, _ := rMeObj.(*miramodels.Me)
	fmt.Printf("You are now logged in, /u/%s\n", rMe.Name)
}

func ExampleReddit_CodeAuth() {
	reddit := mira.Init(mira.Credentials{
		ClientID:     "clientid",
		ClientSecret: "clientsecret",
		UserAgent:    "MIRA CodeAuth Example v0",
		RedirectURL:  "https://example.com/auth", // This must be equal to the value you set in the reddit app config, or you will get an error!
	})

	// Provide all scopes you need to OAuthConfig
	reddit.OAuthConfig.Scopes = []string{"identity", "submit"}

	// Step 1: Generate URL to redirect user to:
	// The state is a unique value you can use to distinguish between user requests. It won't be used by reddit, just returned to your app later (see below).
	state := "UniqueStateToDistinguish"
	// AuthURL is a link to reddit.com, where a user has to confirm the scopes you request, and can choose to accept/decline the auth request.
	AuthURL := reddit.AuthCodeURL(state, reddit.OAuthConfig.Scopes)
	fmt.Printf("Visit this URL to continue the Authentication Flow: %s\n", AuthURL)

	// Step 2: If a user approves your app, he will be redirected to your RedirectURL.
	// In this example, the request would look like https://example.com/auth?state=UniqueStateToDistinguish&code=TOPSECRETCODE
	// With the code, you can then continue with Step 2 of the authentication flow: Actually getting your access token.
	// The code below usually is in a separate function that handles the "auth" endpoint. Just re-initialize "reddit" like above.
	code := "TOPSECRETCODE" // usually would be sth. like r.URL.Query()["code"]
	if err := reddit.CodeAuth(code, nil); err != nil {
		panic(err)
	}

	// Done! You can now use reddit authenticated as that user & can use reddit:
	rMeObj, err := reddit.Me().Info()
	if err != nil {
		panic(err)
	}
	rMe, _ := rMeObj.(*miramodels.Me)
	fmt.Printf("You are now logged in, /u/%s\n", rMe.Name)

	// Note that the access token cannot be retreived separately. This is by design, since it is useless after a few hours.
	// Instead, you should use a TokenNotifyFunc. See example there.
}

// Assumung you have saved the Refresh Token (see example for TokenNotifyFunc), this is how you can restore a session using it:
func ExampleReddit_CodeAuth_resumingSession() {
	handleRefreshToken := func(t *oauth2.Token) error {
		// See TokenNotifyFunc example
		return nil
	}
	reddit := mira.Init(mira.Credentials{
		ClientID:     "clientid",
		ClientSecret: "clientsecret",
		UserAgent:    "MIRA CodeAuth Example v0",
		RedirectURL:  "https://example.com/auth", // This must be equal to the value you set in the reddit app config, or you will get an error!
	})
	reddit.OAuthConfig.Scopes = []string{"identity", "submit"}

	refreshToken := "secret" // fetched from database or similar

	reddit.SetToken(&oauth2.Token{
		AccessToken:  "VOID",
		RefreshToken: refreshToken,
		Expiry:       time.Now().Add(time.Minute * -1), // Set the Access Token to expire immediately so a new one is fetched
	}, reddit.OAuthConfig.Scopes, handleRefreshToken)

	// Done! You can now use reddit authenticated as that user & can use reddit:
	rMeObj, err := reddit.Me().Info()
	if err != nil {
		panic(err)
	}
	rMe, _ := rMeObj.(*miramodels.Me)
	fmt.Printf("You are now logged in again, /u/%s\n", rMe.Name)
}

// This example assumes you have read and understand the CodeAuth() example.
//
// To save the refresh token & grab a new access token once it expires, you can use a notify function that will be called
// each time the token is refreshed (which happens automatically while your app is running)
func ExampleTokenNotifyFunc() {
	reddit := mira.Init(mira.Credentials{
		ClientID:     "clientid",
		ClientSecret: "clientsecret",
		UserAgent:    "MIRA TokenNotifyFunc Example v0",
		RedirectURL:  "https://example.com/auth",
	})
	reddit.OAuthConfig.Scopes = []string{"identity", "submit"}

	// handleRefreshToken statifies the mira.TokenNotifyFunc interface
	handleRefreshToken := func(t *oauth2.Token) error {
		fmt.Println("Refreshed Token:")
		fmt.Println("- New Access Token:", t.AccessToken)
		fmt.Println("- New Refresh Token:", t.RefreshToken)

		// Probably save to Database or similar.
		return nil
	}

	// Step 1 is similar to the CodeAuth example. See there.
	// Step 2 changes slightly. For extended documentation please see the CodeAuth example.
	code := "TOPSECRETCODE" // usually would be sth. like r.URL.Query()["code"]
	if err := reddit.CodeAuth(code, handleRefreshToken); err != nil {
		panic(err)
	}

	// This is it! Each time the token will be refreshed (which happens automatically), handleRefreshToken() will be called with the latest information.
	// If you have multiple users, you should probably add something to distinguish between them :)
}

// If you have multiple users in parallel, you can initialize multiple reddit objects to switch between them:
func ExampleReddit() {
	redditInstances := make(map[string]*mira.Reddit)
	// Login as users - heavily shortened, see CodeAuth / LoginAuth examples:
	// User 1: ttgmpsn
	reddit1 := mira.Init(mira.Credentials{})
	reddit1.CodeAuth("", nil)
	redditInstances["ttgmpsn"] = reddit1

	// User 2: ttgmbot
	reddit2 := mira.Init(mira.Credentials{})
	reddit2.CodeAuth("", nil)
	redditInstances["ttgmbot"] = reddit2

	// If redditInstances is a global variable, you can now use it everywhere!
}

func ExampleReddit_StreamPosts() {
	// Initialize reddit instance like usually - see other examples.
	reddit := mira.Init(mira.Credentials{})

	// Create stream
	stream, err := reddit.Subreddit("pics").StreamPosts()
	if err != nil {
		panic(err)
	}

	// Create listener
	go func() {
		var s miramodels.Submission
		for s = range stream.C {
			if s == nil {
				fmt.Println("Stream was closed")
				return
			}
			fmt.Println("Received new item in stream:", s.GetID())
		}
	}()
}
