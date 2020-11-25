# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## v0.26.0
- [Breaking] Refactored config API. Master data is stored as JSON
- Remove toml dependency

## v0.25.0
- [NEW] HSDP PKI API support

## v0.24.0
- Logging fixes
- [Cartel] Minor API tweaks for better error reporting

## v0.23.0
- Maintenance release
- Add Tag() option for Cartel

## v0.22.0
- [NEW] Console API support: Autoscalers

## v0.21.1
- Fix config URL

## v0.21.0
- [NEW] Hosted appstream (HAS) support
- [NEW] IAM Password policies management

## v0.20.0
- Autoconfig support for IAM, Cartel and Logging
- Add ap3 region
- Fallback mechanism for discovery

## v0.19.0
- [NEW] Service/config discovery
- Share structs with gautocloud-connectors

## v0.18.0
- [NEW] IronIO Worker support
  Codes CRUD
  Tasks CRUD
  Schedules CRUD
  Clusters Read/stats

## v0.17.0
- [IAM] Metadata field changes

## v0.16.0
- [IAM] Add UpdateClient() method 
- [Logging] Add Meta tags to Resource

## v0.15.0

- [NEW] Cartel API support
- [IAM] Switch to SCIM based organization management
- [IAM] Add Organizations.DeleteOrganization()
- [IAM] Add Organizations.DeleteStatus()
- [Logging] Detect errors in batch sends

## v0.14.0
- [IAM] Move user find API to v2
- [IAM] Update go-hsdp-signer
- [Logging] Better support for custom logging

## v0.2.0
- Upgrade github.com/Jeffail/gabs

## v0.1.0
- IAM support
- Logging support
- TDR basic support

[Unreleased]: https://github.com/philips-software/go-hsdp-api/compare/1.0.0...HEAD

