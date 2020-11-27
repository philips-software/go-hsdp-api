# Config
Services have specific endpoints per region. This config API helps you
discover the available services and their settings. Since the data is
machine readable this enables auto configuration of endpoints.

# Canonical source
The canonical source for configuration is:
https://github.com/philips-software/go-hsdp-api/blob/master/config/hsdp.json

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
	baseIAMURLInUSEastClientTest := c.Region("us-east").Env("client-test").Service("iam").URL
	fmt.Printf("IAM Base URL: %s\n", baseIAMURLInUSEastClientTest)
}
```
