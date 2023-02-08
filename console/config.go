package console

import "io"

// Config contains the configuration of a client
type Config struct {
	Region         string
	BaseConsoleURL string
	UAAURL         string
	MetricsAPIURL  string
	Scopes         []string
	Debug          bool
	DebugLog       io.Writer
}
