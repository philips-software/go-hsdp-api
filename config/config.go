//go:generate go-bindata -pkg config -o bindata.go hsdp.toml
package config

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pelletier/go-toml"
)

const (
	// CanonicalURL is the source of truth
	CanonicalURL = "https://raw.githubusercontent.com/philips-software/go-hsdp-api/master/config/hsdp.toml"
)

// Config holds the state of a Config instance
type Config struct {
	region      string
	environment string
	source      io.Reader
	config      *toml.Tree
}

// Service holds the relevant toml data for a service
type Service struct {
	config *toml.Tree
}

type OptionFunc func(*Config) error

// New returns a Config Instance. You can pass
// a list OptionFunc to cater the Config to your needs
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
			// Fallback to baked in copy in case github.com is down,
			// but only if its not older than 180 days
			if bakedInCopy, err := hsdpToml(); err != nil ||
				bakedInCopy.info.ModTime().Before(time.Now().AddDate(0, 0, -180)) {
				return nil, ErrUnreachableOrOutdatedConfigSource
			} else {
				data, err := Asset(bakedInCopy.info.Name())
				if err != nil {
					return nil, err
					// Asset was not found.
				}
				config.source = bytes.NewReader(data)
			}
		} else {
			defer resp.Body.Close()
			config.source = resp.Body
		}
	}
	content, err := toml.LoadReader(config.source)
	if err != nil {
		return nil, err
	}
	config.config = content
	return config, nil
}

// FromReader option specifies the toml source to read
// If this option is not provided the canonical source
// hosted on Github will be used. See CanonicalURL
func FromReader(reader io.Reader) OptionFunc {
	return func(c *Config) error {
		c.source = reader
		return nil
	}
}

// WithRegion sets the region of the newly created Config instance
func WithRegion(region string) OptionFunc {
	return func(c *Config) error {
		c.region = region
		return nil
	}
}

// WithEnv sets the environment of the newly created Config instance
func WithEnv(env string) OptionFunc {
	return func(c *Config) error {
		c.environment = env
		return nil
	}
}

// Region returns a new Config instance with region set
func (c *Config) Region(region string) *Config {
	return &Config{
		config:      c.config,
		region:      region,
		environment: c.environment,
	}
}

// Env returns a new Config instance with environment set
func (c *Config) Env(environment string) *Config {
	return &Config{
		config:      c.config,
		region:      c.region,
		environment: environment,
	}
}

// Services returns a list of available services in the region
func (c *Config) Services() []string {
	services := make([]string, 0)
	// region level
	regional, ok := c.config.Get(fmt.Sprintf("region.%s.service", c.region)).(*toml.Tree)
	if ok && len(regional.Keys()) > 0 {
		services = append(services, regional.Keys()...)
	}
	// environment
	environment, ok := c.config.Get(fmt.Sprintf("region.%s.env.%s.service", c.region, c.environment)).(*toml.Tree)
	if ok && len(environment.Keys()) > 0 {
		services = append(services, environment.Keys()...)
	}
	return services
}

// Service returns an instance scoped to the service in the region and environment
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

// String returns the string key of the Service
func (s *Service) String(key string) (string, error) {
	if s.config == nil {
		return "", ErrMissingConfig
	}
	out, ok := s.config.Get(key).(string)
	if ok {
		return out, nil
	}
	return "", ErrNotFound
}

// Int returns the int key of the Service
func (s *Service) Int(key string) (int, error) {
	if s.config == nil {
		return 0, ErrMissingConfig
	}
	out, ok := s.config.Get(key).(int)
	if ok {
		return out, nil
	}
	return 0, ErrNotFound
}

// Available returns true if the Service exists and has data
func (s *Service) Available() bool {
	if s == nil || s.config == nil {
		return false
	}
	return len(s.config.Keys()) != 0
}

// Keys returns the list of available keys for Service
func (s *Service) Keys() []string {
	keys := make([]string, 0)
	if s == nil || s.config == nil {
		return keys
	}
	keys = append(keys, s.config.Keys()...)
	return keys
}
