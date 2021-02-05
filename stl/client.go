package stl

import (
	"context"
	"github.com/hasura/go-graphql-client"
	"github.com/philips-software/go-hsdp-api/console"
	"github.com/philips-software/go-hsdp-api/internal"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

const (
	userAgent = "go-hsdp-api/stl/" + internal.LibraryVersion
)

// OptionFunc is the function signature function for options
type OptionFunc func(*http.Request) error

// Config contains the configuration of a client
type Config struct {
	Region      string
	Environment string
	STLAPIURL   string
	DebugLog    string
}

// A Client manages communication with HSDP DICOM API
type Client struct {
	// HTTP client used to communicate with IAM API
	consoleClient *console.Client

	gql *graphql.Client

	config *Config

	// User agent used when communicating with the HSDP DICOM API.
	UserAgent string

	debugFile *os.File

	Devices   *DevicesService
	Resources *ResourcesService
}

// NewClient returns a new HSDP DICOM API client. Configured console and IAM clients
// must be provided as the underlying API requires tokens from respective services
func NewClient(consoleClient *console.Client, config *Config) (*Client, error) {
	return newClient(consoleClient, config)
}

func newClient(consoleClient *console.Client, config *Config) (*Client, error) {
	c := &Client{consoleClient: consoleClient, config: config, UserAgent: userAgent}
	httpClient := oauth2.NewClient(context.Background(), consoleClient)

	if config.DebugLog != "" {
		var err error
		c.debugFile, err = os.OpenFile(config.DebugLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err == nil {
			httpClient.Transport = newLoggingRoundTripper(httpClient.Transport, c.debugFile)
		}
	}
	header := make(http.Header)
	header.Set("User-Agent", userAgent)
	httpClient.Transport = newHeaderRoundTripper(httpClient.Transport, header)

	c.gql = graphql.NewClient("https://console.na3.hsdp.io/api/stl/user/v1/graphql", httpClient)
	c.Devices = &DevicesService{client: c}
	c.Resources = &ResourcesService{client: c}

	return c, nil
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
