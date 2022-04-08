package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type UserConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type WebsiteConfig struct {
	Domain string `yaml:"domain" json:"domain"`
	Title  string `yaml:"title" json:"title"`
}

type TinylyticsConfig struct {
	User       UserConfig      `yaml:"user"`
	Websites   []WebsiteConfig `yaml:"websites"`
	DataFolder string          `yaml:"data-folder"`
}

var Config TinylyticsConfig

func LoadConfig(path string) {
	err := cleanenv.ReadConfig(path, &Config)
	if err != nil {
		fmt.Println("Error loading config", err)
	}
}
