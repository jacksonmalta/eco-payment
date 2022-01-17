package authorizer

type Config struct {
	Url string
}

func (c *Config) WithUrl(url string) *Config {
	c.Url = url
	return c
}
