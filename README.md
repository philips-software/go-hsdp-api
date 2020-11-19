[![Build Status](https://travis-ci.com/philips-software/go-hsdp-api.svg?branch=master)](https://travis-ci.com/philips-software/go-hsdp-api)
[![Maintainability](https://api.codeclimate.com/v1/badges/125caa4282d4d82b84cd/maintainability)](https://codeclimate.com/github/philips-software/go-hsdp-api/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/125caa4282d4d82b84cd/test_coverage)](https://codeclimate.com/github/philips-software/go-hsdp-api/test_coverage)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/philips-software/go-hsdp-api)](https://pkg.go.dev/github.com/philips-software/go-hsdp-api)

# go-hsdp-api

A HSDP API client enabling Go programs to interact with various HSDP APIs in a simple and uniform way

## Supported APIs

The current implement covers only a subset of HSDP APIs. Basically we implement functionality as needed.

- [x] Cartel c.q. Container Host management ([examples](cartel/README.md))
- [x] IronIO tasks, codes and schedules management ([examples](iron/README.md))
- [x] HSDP PKI services
- [x] HSDP IAM/IDM management
  - [x] Groups
  - [x] Organizations
  - [x] Permissions
  - [x] Roles
  - [x] Users
  - [x] Passwords
  - [x] Propositions
  - [x] Applications
  - [x] Services
  - [x] Devices
  - [x] MFA Policies
  - [x] Password Policies
- [x] Logging ([examples](logging/README.md))
- [ ] Auditing
- [x] S3 Credentials Policy management
- [x] Hosted Application Streaming (HAS) management
- [x] Configuration
- [x] Console settings
  - [ ] Metrics Alerts
  - [x] Metrics Autoscalers

## Usage

```go
package main

import (
        "fmt"

        "github.com/philips-software/go-hsdp-api/iam"
)

func main() {
        client, _ := iam.NewClient(nil, &iam.Config{
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
        introspect, _, _ := client.Introspect()
        if introspect != nil {
                fmt.Printf("Introspect response: %v\n", introspect)
        }
}
```

## TODO

- Increase API coverage

## Issues

- If you have an issue: report it on the [issue tracker](https://github.com/philips-software/go-hsdp-api/issues)

## Contact / Getting help

Andy Lo-A-Foe (<andy.lo-a-foe@philips.com>)

## License

License is MIT. See [LICENSE file](LICENSE.md)
