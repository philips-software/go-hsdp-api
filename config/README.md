# Config
Many services have specific endpoints per region. Th config API helps you
discover the available services and their settings

# Canonical source
The canonical source for configuration is:
https://github.com/philips-software/go-hsdp-api/blob/master/config/hsdp.toml

# Example
```go
package main

import (
	"fmt"

	"github.com/philips-software/go-hsdp-api/config"
)

func main() {
	c, err := config.New()
	if err != nil {
		fmt.Printf("error loading: %v\n", err)
		return
	}
	baseIAMURLInUSEastClientTest, err := c.Region("us-east").Env("client-test").Service("iam").String("url")
	if err != nil {
		fmt.Printf("not found: %v\n", err)
		return
	}
	fmt.Printf("IAM Base URL: %s\n", baseIAMURLInUSEastClientTest)
```
