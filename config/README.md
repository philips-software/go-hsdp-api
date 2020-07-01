# Config
Services have specific endpoints per region. Thisconfig API helps you
discover the available services and their settings. Since the data is
machine readable this enables auto configuration of endpoints.

# Canonical source
The canonical source for configuration is:
https://github.com/philips-software/go-hsdp-api/blob/master/config/hsdp.toml

# Example
Determine the IAM base URL of a region and environment

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
	baseIAMURLInUSEastClientTest, err := c.Region("us-east").Env("client-test").Service("iam").String("iam_url")
	if err != nil {
		fmt.Printf("not found: %v\n", err)
		return
	}
	fmt.Printf("IAM Base URL: %s\n", baseIAMURLInUSEastClientTest)
}
```
