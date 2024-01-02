package container

import (
	"bahamut/internal/config"
	"net/http"
)

type Container struct {
	httpClient *http.Client
	config     *config.Config
}

func Register(h *http.Client, c *config.Config) *Container {
	return &Container{
		httpClient: h,
		config:     c,
	}
}

func (c *Container) HttpClient() *http.Client {
	return c.httpClient
}

func (c *Container) Config() *config.Config {
	return c.config
}
