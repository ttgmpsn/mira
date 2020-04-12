package mira

// Init will initialize the Reddit instance.
// Note that you most likely want to auth using
// LoginAuth() or CodeAuth()
func Init(c Credentials) *Reddit {
	auth := newOAuthSession(c)
	auth.SetDefault()
	return auth
}
