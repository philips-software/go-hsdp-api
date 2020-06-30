package config

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pelletier/go-toml"
)

const (
	CanonicalURL = "https://raw.githubusercontent.com/philips-software/go-hsdp-api/master/config/hsdp.toml"
)

type Config struct {
	region      string
	environment string
	source      io.Reader
	config      *toml.Tree
}

type Service struct {
	config *toml.Tree
}

type OptionFunc func(*Config) error

// New returns a Config Instance
func New(opts ...OptionFunc) (*Config, error) {
	config := &Config{}
	for _, opt := range opts {
		if err := opt(config); err != nil {
			return nil, err
		}
	}
	if config.source == nil {
		resp, err := http.Get(CanonicalURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		config.source = resp.Body
	}
	content, err := toml.LoadReader(config.source)
	if err != nil {
		return nil, err
	}
	config.config = content
	return config, nil
}

func FromReader(reader io.Reader) OptionFunc {
	return func(c *Config) error {
		c.source = reader
		return nil
	}
}

func WithRegion(region string) OptionFunc {
	return func(c *Config) error {
		c.region = region
		return nil
	}
}

func WithEnv(env string) OptionFunc {
	return func(c *Config) error {
		c.environment = env
		return nil
	}
}

func (c *Config) Region(region string) *Config {
	return &Config{
		config:      c.config,
		region:      region,
		environment: c.environment,
	}
}

func (c *Config) Env(environment string) *Config {
	return &Config{
		config:      c.config,
		region:      c.region,
		environment: environment,
	}
}

func (c *Config) Service(service string) *Service {
	// Check if service is at region level
	regionService, ok := c.config.Get(fmt.Sprintf("region.%s.service.%s", c.region, service)).(*toml.Tree)
	if ok {
		return &Service{
			config: regionService,
		}
	}
	// Otherwise check at environment level
	envService, ok := c.config.Get(fmt.Sprintf("region.%s.env.%s.service.%s", c.region, c.environment, service)).(*toml.Tree)
	if ok {
		return &Service{
			config: envService,
		}
	}
	return &Service{}
}

func (s *Service) String(str string) (string, error) {
	if s.config == nil {
		return "", fmt.Errorf("missing config")
	}
	out, ok := s.config.Get(str).(string)
	if ok {
		return out, nil
	}
	return "", fmt.Errorf("not found")
}

func (s *Service) Available() bool {
	if s == nil || s.config == nil {
		return false
	}
	return len(s.config.Keys()) != 0
}
