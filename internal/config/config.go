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
	Login      *ConfigModulesLogin
	Sign       *ConfigModulesSign
	LineNotify *ConfigModulesLineNotify `yaml:"line_notify"`
}

type ConfigModulesLogin struct {
	Username string
	Password string
	Debug    bool
}

type ConfigModulesSign struct {
	Debug bool
}

type ConfigModulesLineNotify struct {
	Token string `yaml:"personal_token"`
}

type ConfigBrowser struct {
	Driver           string   `yaml:"driver"`
	Args             []string `yaml:"args"`
	Headless         bool     `yaml:"headless"`
	SkipInstallation bool     `yaml:"skip_installation"`
	Debug            bool     `yaml:"debug"`
}

func InitConfig(filePath string) (*Config, error) {
	var config Config
	fileBuffer, err := os.ReadFile(filePath)
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
