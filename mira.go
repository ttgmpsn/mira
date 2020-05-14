package mira

// Init will initialize the Reddit instance.
// Note that you most likely want to auth using
// LoginAuth() or CodeAuth()
func Init(c Credentials) *Reddit {
	instance := newOAuthSession(c)
	instance.Chain = make(chan *chainVals, 32)
	instance.SetDefault()
	return instance
}
