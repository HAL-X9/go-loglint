package config

type Rule struct {
	Name   string `yaml:"name"`
	Enable bool   `yaml:"enable"`
}

type Loggers struct {
	Include []string `yaml:"include"`
}

type Config struct {
	AutoFix bool    `yaml:"auto_fix"`
	Rules   []Rule  `yaml:"rules"`
	Loggers Loggers `yaml:"loggers"`
}
