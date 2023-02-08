// Package stl provides support for HSDP STL services
package stl

import (
	"context"
	"io"
	"net/http"

	"github.com/hasura/go-graphql-client"
	autoconf "github.com/philips-software/go-hsdp-api/config"
	"github.com/philips-software/go-hsdp-api/console"
	"github.com/philips-software/go-hsdp-api/internal"
	"golang.org/x/oauth2"
)

const (
	userAgent = "go-hsdp-api/edge/" + internal.LibraryVersion
)

// OptionFunc is the function signature function for options
type OptionFunc func(*http.Request) error

// Config contains the configuration of a consoleClient
type Config struct {
	Region      string
	Environment string
	STLAPIURL   string
	DebugLog    io.Writer
}

// A Client manages communication with HSDP Edge API
type Client struct {
	// HTTP consoleClient used to communicate with IAM API
	consoleClient *console.Client

	gql *graphql.Client

	config *Config

	// User agent used when communicating with the HSDP Edge API.
	UserAgent string

	Devices *DevicesService
	Apps    *AppsService
	Config  *ConfigService
	Certs   *CertsService
}

// NewClient returns a new HSDP Edge API consoleClient. Configured console and IAM clients
// must be provided as the underlying API requires tokens from respective services
func NewClient(consoleClient *console.Client, config *Config) (*Client, error) {
	return newClient(consoleClient, config)
}

func newClient(consoleClient *console.Client, config *Config) (*Client, error) {
	doAutoconf(config)
	c := &Client{consoleClient: consoleClient, config: config, UserAgent: userAgent}
	httpClient := oauth2.NewClient(context.Background(), consoleClient)

	if config.DebugLog != nil {
		httpClient.Transport = internal.NewLoggingRoundTripper(httpClient.Transport, config.DebugLog)
	}
	header := make(http.Header)
	header.Set("User-Agent", userAgent)
	httpClient.Transport = internal.NewHeaderRoundTripper(httpClient.Transport, header)

	c.gql = graphql.NewClient(config.STLAPIURL, httpClient)
	c.Devices = &DevicesService{client: c}
	c.Apps = &AppsService{client: c}
	c.Config = &ConfigService{client: c}
	c.Certs = &CertsService{client: c}

	return c, nil
}

func doAutoconf(config *Config) {
	if config.Region != "" {
		c, err := autoconf.New(
			autoconf.WithRegion(config.Region))
		if err == nil {
			stlService := c.Service("stl")
			if config.STLAPIURL == "" {
				config.STLAPIURL = stlService.URL
			}
		}
	}
}

// Query is a generic GraphQL query
func (c *Client) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return c.gql.Query(ctx, q, variables)
}

// Close releases allocated resources of clients
func (c *Client) Close() {
}
