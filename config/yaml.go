package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database       string `yaml:"database"`
	SecreteKey     string `yaml:"secret_key"`
	SecretKeyReset string `yaml:"secret_key_reset"`
}

func NewConfig() *Config {
	cpath := "config.yml"
	var config Config

	// Open config file
	file, err := os.Open(cpath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		panic(err)
	}

	return &config
}
