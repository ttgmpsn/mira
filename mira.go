package mira

// Init will initialize the Reddit instance.
// Note that you most likely want to auth using
// LoginAuth() or CodeAuth()
func Init(c Credentials) *Reddit {
	auth := NewOAuthSession(c)
	auth.SetDefault()
	return auth
}
