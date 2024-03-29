package scraper

type Config struct {
	Query Query `yaml:"query"`
}

func (c *Config) Prefix() string {
	return "ScraperConfig"
}

type Query struct {
	StartIndex int `yaml:"startIndex"`
	Limit      int `yaml:"limit"`
}
