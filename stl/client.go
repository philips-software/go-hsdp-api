package stl

import (
	"github.com/philips-software/go-hsdp-api/console"
	"github.com/philips-software/go-hsdp-api/internal"
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

	config *Config

	// User agent used when communicating with the HSDP DICOM API.
	UserAgent string

	debugFile *os.File

	Devices *DevicesService
}

// NewClient returns a new HSDP DICOM API client. Configured console and IAM clients
// must be provided as the underlying API requires tokens from respective services
func NewClient(consoleClient *console.Client, config *Config) (*Client, error) {
	return newClient(consoleClient, config)
}

func newClient(consoleClient *console.Client, config *Config) (*Client, error) {
	c := &Client{consoleClient: consoleClient, config: config, UserAgent: userAgent}
	if config.DebugLog != "" {
		var err error
		c.debugFile, err = os.OpenFile(config.DebugLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			c.debugFile = nil
		}
	}
	c.Devices = &DevicesService{client: c}

	return c, nil
}

// Close releases allocated resources of clients
func (c *Client) Close() {
	if c.debugFile != nil {
		_ = c.debugFile.Close()
		c.debugFile = nil
	}
}
