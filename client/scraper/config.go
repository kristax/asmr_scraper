package scraper

type Config struct {
	ForceUpdateInfo  bool `yaml:"forceUpdateInfo"`
	ForceUploadImage bool `yaml:"forceUploadImage"`
}

func (c *Config) Prefix() string {
	return "ScraperConfig"
}
