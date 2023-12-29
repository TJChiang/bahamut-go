package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Modules *ConfigModules
	Browser *ConfigBrowser
}

type ConfigModules struct {
	Login *ConfigModulesLogin
	Sign  *ConfigModulesSign
}

type ConfigModulesLogin struct {
	Username string
	Password string
	Debug    bool
}

type ConfigModulesSign struct {
	Debug bool
}

type ConfigBrowser struct {
	Driver           string   `yaml:"driver"`
	Type             string   `yaml:"type"`
	Args             []string `yaml:"args"`
	Headless         bool     `yaml:"headless"`
	SkipInstallation bool     `yaml:"skip_installation"`
	Debug            bool     `yaml:"debug"`
}

func InitConfig() (*Config, error) {
	var config Config
	fileBuffer, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(fileBuffer), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) GetModulesConfig() *ConfigModules {
	return c.Modules
}

func (c *Config) GetBrowserConfig() *ConfigBrowser {
	return c.Browser
}
