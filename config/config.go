// Package world provides info on HSDP services
package config

import (
	"embed"
	_ "embed"
	"encoding/json"
	"io"
)

// Config holds the state of a Config instance
type Config struct {
	region      string
	environment string
	source      io.Reader
	world       World
}

type World struct {
	Regions map[string]Region `json:"region"`
}
type Region struct {
	Environments map[string]Environment `json:"env,omitempty"`
	Services     map[string]Service     `json:"service,omitempty"`
}

type Environment struct {
	Services map[string]Service `json:"service,omitempty"`
}

// Service holds the relevant data for a service
type Service struct {
	URL    string `json:"url,omitempty"`
	Domain string `json:"domain,omitempty"`
	Host   string `json:"host,omitempty"`
}

type OptionFunc func(*Config) error

//go:embed hsdp.json
var cfg embed.FS

// New returns a Config Instance. You can pass
// a list OptionFunc to cater the Config to your needs
func New(opts ...OptionFunc) (*Config, error) {
	config := &Config{}
	for _, opt := range opts {
		if err := opt(config); err != nil {
			return nil, err
		}
	}
	var world World
	data, err := cfg.ReadFile("hsdp.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &world)
	if err != nil {
		return nil, err
	}
	config.world = world
	return config, nil
}

// FromReader option specifies the JSON source to read
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
		c.region = c.regionMapping(region)
		return nil
	}
}

// WithEnv sets the environment of the newly created Config instance
func WithEnv(env string) OptionFunc {
	if env == "production" {
		env = "prod"
	}
	return func(c *Config) error {
		c.environment = env
		return nil
	}
}

// Region returns a new Config instance with region set
func (c *Config) Region(region string) *Config {
	return &Config{
		world:       c.world,
		region:      c.regionMapping(region),
		environment: c.environment,
	}
}

// Env returns a new Config instance with environment set
func (c *Config) Env(environment string) *Config {
	return &Config{
		world:       c.world,
		region:      c.region,
		environment: environment,
	}
}

// Regions returns the known regions
func (c *Config) Regions() []string {
	regions := make([]string, 0)
	// region level
	for k := range c.world.Regions {
		regions = append(regions, k)
	}
	return regions
}

// Services returns a list of available services in the region
func (c *Config) Services() []string {
	services := make([]string, 0)
	// region level
	if svcs, ok := c.world.Regions[c.region]; ok {
		for s := range svcs.Services {
			services = append(services, s)
		}
	}

	// environment
	if svcs, ok := c.world.Regions[c.region].Environments[c.environment]; ok {
		for s := range svcs.Services {
			services = append(services, s)
		}
	}
	return services
}

// Service returns an instance scoped to the service in the region and environment
func (c *Config) Service(service string) *Service {
	// Check if service is at region level
	if regionService, ok := c.world.Regions[c.region]; ok {
		if service, ok := regionService.Services[service]; ok {
			return &service
		}
		if envService, ok := regionService.Environments[c.environment]; ok {
			if service, ok := envService.Services[service]; ok {
				return &service
			}
		}
	}
	return &Service{}
}

func (c *Config) regionMapping(region string) string {
	switch region {
	case "us-east-1":
		return "us-east"
	case "eu-west-1":
		return "eu-west"
	case "ap-se-2":
		return "apac3"
	case "sa-east-1":
		return "sa1"
	default:
		return region
	}
}
