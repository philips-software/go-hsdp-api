[![Build Status](https://travis-ci.com/philips-software/go-hsdp-api.svg?branch=master)](https://travis-ci.com/philips-software/go-hsdp-api)
[![Maintainability](https://api.codeclimate.com/v1/badges/125caa4282d4d82b84cd/maintainability)](https://codeclimate.com/github/philips-software/go-hsdp-api/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/125caa4282d4d82b84cd/test_coverage)](https://codeclimate.com/github/philips-software/go-hsdp-api/test_coverage)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/philips-software/go-hsdp-api)](https://pkg.go.dev/github.com/philips-software/go-hsdp-api)

# go-hsdp-api

A HSDP API client enabling Go programs to interact with various HSDP APIs in a simple and uniform way

## Supported APIs

The current implement covers only a subset of HSDP APIs. Basically we implement functionality as needed.

- [x] Cartel c.q. Container Host management ([examples](cartel/README.md))
- [x] Secure Transport Layer (STL) / Edge 
  - [x] Device queries
  - [x] Application Resources management
  - [x] Device configuration management (firewall, logging)
- [x] Public Key Infrastructure (PKI) management
- [x] Identity and Access Management (IAM)
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
  - [x] Email Templates
  - [x] SMS Gateways
  - [ ] SMS Templates
- [x] Logging ([examples](logging/README.md))
- [x] Auditing ([examples](audit/README.md))
- [x] Clinical Data Repository (CDR)
  - [x] Tenant Onboarding
  - [x] Subscription management
  - [x] FHIR CRUD
  - [x] FHIR Patch
  - [x] STU3
  - [x] R4
- [x] Telemetry Data Repository (TDR)
  - [x] Contract management
  - [x] Data Item management
- [x] S3Creds Policy management
- [x] DICOM Store
  - [x] Config management
- [x] Notification service
- [x] Hosted Application Streaming (HAS) management
- [x] Service Discovery
- [x] Console settings
  - [ ] Metrics Alerts
  - [x] Metrics Autoscalers
- [x] IronIO tasks, codes and schedules management ([examples](iron/README.md))
- [x] Clinical Data Lake (CDL) management
  - [x] Research Studies
  - [x] Data Type Definitions
  - [x] Label Definitions
  - [x] Export Routes
- [x] AI Inference
  - [x] Compute Environment management
  - [x] Compute Target managements
  - [x] Model management
  - [x] Inference Job management
- [x] AI Training
  - [x] Compute Environment management
  - [x] Model management
- [x] AI Workspace
  - [x] Compute Target management
  - [x] Workspace management

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
