// Package docker provides support for HSDP Docker Registry services
package docker

import (
	"context"
	"net/http"
	"os"

	"github.com/hasura/go-graphql-client"
	autoconf "github.com/philips-software/go-hsdp-api/config"
	"github.com/philips-software/go-hsdp-api/console"
	"github.com/philips-software/go-hsdp-api/internal"
)

const (
	userAgent = "go-hsdp-api/docker/" + internal.LibraryVersion
)

// OptionFunc is the function signature function for options
type OptionFunc func(*http.Request) error

// Config contains the configuration of a consoleClient
type Config struct {
	Region       string
	DockerAPIURL string
	DebugLog     string
	host         string
}

// A Client manages communication with HSDP DICOM API
type Client struct {
	// HTTP consoleClient used to communicate with IAM API
	*console.Client

	gql *graphql.Client

	config *Config

	// User agent used when communicating with the HSDP DICOM API.
	UserAgent string

	debugFile *os.File

	ServiceKeys  *ServiceKeysService
	Namespaces   *NamespacesService
	Repositories *RepositoriesService
}

// NewClient returns a new HSDP Docker Registry API client. A configured console client
// must be provided as the underlying API requires tokens from respective service
func NewClient(consoleClient *console.Client, config *Config) (*Client, error) {
	return newClient(consoleClient, config)
}

func newClient(consoleClient *console.Client, config *Config) (*Client, error) {
	doAutoconf(config)

	c := &Client{Client: consoleClient, config: config, UserAgent: userAgent}

	header := make(http.Header)
	header.Set("Accept", "*/*")

	// Injecting these headers so we satisfy the proxies
	consoleClient.Transport = internal.NewHeaderRoundTripper(consoleClient.Transport, header, func(req *http.Request) error {
		token, err := consoleClient.Token()
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "bearer "+token.AccessToken)
		req.Header.Set("X-User-Access-Token", token.AccessToken)
		return nil
	})

	c.gql = graphql.NewClient(config.DockerAPIURL, consoleClient.Client)
	c.ServiceKeys = &ServiceKeysService{client: c}
	c.Namespaces = &NamespacesService{client: c}
	c.Repositories = &RepositoriesService{client: c}

	return c, nil
}

func doAutoconf(config *Config) {
	if config.Region != "" {
		c, err := autoconf.New(
			autoconf.WithRegion(config.Region))
		if err == nil {
			dockerService := c.Service("docker-registry")
			if config.DockerAPIURL == "" {
				config.DockerAPIURL = dockerService.URL
			}
			config.host = dockerService.Host
		}
	}
}

// Host returns the regional host base
func (c *Client) Host() string {
	return c.config.host
}

// Query is a generic GraphQL query
func (c *Client) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return c.gql.Query(ctx, q, variables)
}

// Close releases allocated resources of clients
func (c *Client) Close() {
	if c.debugFile != nil {
		_ = c.debugFile.Close()
		c.debugFile = nil
	}
}
