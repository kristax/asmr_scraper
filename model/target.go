package model

type Target struct {
	Id      string `yaml:"id"`
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Async   bool   `yaml:"async"`
	Jitter  int    `yaml:"jitter"`
	Disable bool   `yaml:"disable"`
}
