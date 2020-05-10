package config

type Migration struct {
	Directory string `yaml:"directory"`
	Driver    string `yaml:"driver"`
	DSN       string `yaml:"dsn"`
}
