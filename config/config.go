package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Database struct {
	User   string `yaml:"user"`
	Pwd    string `yaml:"pwd"`
	DbName string `yaml:"db_name"`
	Host   string `yaml:"host"`
}

func (db *Database) CreateDsn() string {
	return fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s sslmode=disable",
		db.Host,
		db.User,
		db.DbName,
		db.Pwd,
	)
}

type CaptchaConfig struct {
	Active    bool   `yaml:"active"`
	SecretKey string `yaml:"secret_key"`
}

type ConfigData struct {
	Database       Database      `yaml:"database"`
	SecreteKey     string        `yaml:"secret_key"`
	SecretKeyReset string        `yaml:"secret_key_reset"`
	MailUrl        string        `yaml:"mail_url"`
	DevMode        bool          `yaml:"dev_mode"`
	Captcha        CaptchaConfig `yaml:"captcha"`
}

var Config *ConfigData

func init() {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Warning: %s environment variable not set.", k)
		}
		return v
	}

	devconfig := os.Getenv("devconfig")

	if devconfig != "" {
		Config = LoadConfig(devconfig)
		Config.DevMode = true
	} else {
		captchaactive := mustGetenv("CAPTCHA_ACTIVE") == "true"
		Config = &ConfigData{
			Database: Database{
				User:   mustGetenv("DB_USER"),
				Pwd:    mustGetenv("DB_PWD"),
				DbName: mustGetenv("DB_NAME"),
				Host:   mustGetenv("DB_HOST"),
			},
			Captcha: CaptchaConfig{
				Active:    captchaactive,
				SecretKey: mustGetenv("CAPTCHA_SECRET_KEY"),
			},
			SecreteKey:     mustGetenv("SECRET_KEY"),
			SecretKeyReset: mustGetenv("SECRET_KEY_RESET"),
		}
		Config.DevMode = false
	}
}

func LoadConfig(cpath string) *ConfigData {
	var config ConfigData

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
