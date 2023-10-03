package asmr_one

type Config struct {
	Host  string `yaml:"host"`
	Debug bool   `yaml:"debug"`
}

func (c *Config) Prefix() string {
	return "AsmrOneConfig"
}
