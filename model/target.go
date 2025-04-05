package model

type Target struct {
	Name     string    `yaml:"name"`
	Type     string    `yaml:"type"`
	Async    bool      `yaml:"async"`
	Jitter   int       `yaml:"jitter"`
	Disable  bool      `yaml:"disable"`
	SubItems *SubItems `yaml:"subItems"`
}

type SubItems struct {
	SortBy string `yaml:"sortBy"`
	Fields string `yaml:"fields"`
}
