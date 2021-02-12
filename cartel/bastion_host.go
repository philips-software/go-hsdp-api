package cartel

import (
	autoconf "github.com/philips-software/go-hsdp-api/config"
)

// BastionHost returns the bastion host if it can be guessed correctly
// An empty string is returned otherwise
func (c *Client) BastionHost() string {
	host := c.config.Host
	if host == "" {
		return ""
	}
	region := c.config.Region
	if region == "" {
		region = "us-east"
	}
	ac, err := autoconf.New(autoconf.WithRegion(region))
	if err != nil {
		return ""
	}
	// Search for matching gateway service
	for _, region := range ac.Regions() {
		if ac.Region(region).Service("cartel").Host == host {
			return ac.Region(region).Service("gateway").Host
		}
	}
	return ""
}
