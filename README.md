# go-hsdp-api

A HSDP API client enabling Go programs to interact with various HSDP APIs in a simple and uniform way

## Supported APIs

The current implement covers only a subset of HSDP APIs. Basically we implement functonality as needed.

- [x] IAM token authorization
- [x] Group management
- [x] Organization management
- [x] Permission management
- [x] Role managemnet
- [x] User management
- [x] Password management
- [ ] Device management
- [ ] Policy management
- [x] Proposition management
- [x] Application management
- [x] Logging
- [ ] Auditing

## Usage

```go
package main

import (
        "fmt"

        "github.com/hsdp/go-hsdp-api/iam"
)

func main() {
        client, _ := api.NewClient(nil, &api.Config{
                OAuth2ClientID: "ClientID",
                OAuth2Secret:   "ClientPWD",
                SharedKey:      "KeyHere",
                SecretKey:      "SecretHere",
                IAMURL:         "https://iam-stage.foo-bar.com",
                IDMURL:         "https://idm-stage.foo-bar.com",
        })
        err := client.Login("iam.login@aemian.com", "Password!@#")
        if err != nil {
                fmt.Printf("Error logging in: %v\n", err)
                return
        }
        introspect, resp, _ := client.Introspect()
        if val != nil {
                fmt.Printf("Introspect response: %v\n", introspect)
        }
}
```

## Todo

- Increase API coverage
- Write tests

## Issues

- If you have an issue: report it on the [issue tracker](https://github.com/hsdp/go-hsdp-api/issues)

## Author

Andy Lo-A-Foe (<andy.loafoe@aemian.com>)

## License

Licensed under the Apache License, Version 2.0 (the "License")
