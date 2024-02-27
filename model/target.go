package model

type Target struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Async   bool   `yaml:"async"`
	Jitter  int    `yaml:"jitter"`
	Disable bool   `yaml:"disable"`
}
