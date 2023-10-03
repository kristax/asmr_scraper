package jellyfin

type Config struct {
	Host   string `yaml:"host"`
	ApiKey string `yaml:"apiKey"`
	UserId string `yaml:"userId"`
	Debug  bool   `yaml:"debug"`
}

func (c *Config) Prefix() string {
	return "JellyfinConfig"
}
