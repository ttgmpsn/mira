package mira

// Init will initialize the Reddit instance.
// Note that you most likely want to auth using
// LoginAuth() or CodeAuth() afterwards, see the examples there.
func Init(c Credentials) *Reddit {
	instance := newOAuthSession(c)
	instance.chain = make(chan *chainVals, 32)
	instance.SetDefault()
	return instance
}

// SetDefault gets sensible default values for streams.
func (c *Reddit) SetDefault() {
	c.Config = redditConfig{
		CommentStreamInterval: 45,
		PostStreamInterval:    45,
	}
}
