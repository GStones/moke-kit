package logging

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Type string
}

func LoadConfig() (config Config, err error) {
	err = envconfig.Process("log", &config)
	return
}
