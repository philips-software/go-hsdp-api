# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).
## v0.38.0
- Add DecryptPayload to decrypt Iron payloads
- Export some convenience functions

## v0.37.0
- NEW: support for updating IAM service certificates
- Prepare for CDR changes
- Improve debug log output

## v0.36.6
- Add cn1 endpoints for IAM/IDM

## v0.36.5
- Add additional S3Creds regions

## v0.36.4
- Add all other possible actions

## v0.36.3
- Add ALL_BUCKET to S3 policy action list

## v0.36.2
- Dependency upgrade

## v0.36.1
- Tweak JSON structs

## v0.36.0
- NEW: Support immediate IAM account activation with optional password input
- Support user updates

## v0.35.6
- Wrap errors for better contextual errors

## v0.35.5
- Improved error checking for application and proposition creation

## v0.35.4
- Remove duplicate DICOM logging

## v0.35.3
- Expose TokenRefresh() in DICOM client

## v0.35.2
- Fix STL cert update call issue

## v0.35.1
- Add iam.Applications.GetApplicationByName()

## v0.35.0
- NEW: Secure Transport Layer (STL) support

## v0.34.4
- Fix some PKI methods
- Add IAM token revoke calls

## v0.34.3
- Remove elastic due to license change

## v0.34.2
- Add pki.Services.GetCertificates() method

## v0.34.1
- Fix minimum version in go.mod

## v0.34.0
- Remove online config refresh code
- Use Go 1.16 embed to bundle hsdp.json
- This version only works with Go 1.16+ (breaking change!)

## v0.33.0
- NEW: DICOM Config API support
- Create internal package for consistent better versioning
- S3Creds renaming (breaking change!)

## v0.32.3
- Disable keep-alive for Cartel

## v0.32.2
- Better error reporting in Cartel

## v0.32.1
- Add cartel.BastionHost() 

## v0.32.0
- Proxy support

## v0.31.2
- Email template: change locale handling

## v0.31.1
- Bug fixes

## v0.31.0
- [NEW] IDM Email Templates support

## v0.30.0
- Remove fhir package

## v0.29.1
- Fix HSDP Audit documentation
- Unexport many methods

## v0.29.0
- [NEW] HSDP Audit support

## v0.28.0
- Logging: support bearer token auth (client credentials / service identities)
- Tweaking of CDR client

## v0.27.2
- CDR: Add endpoint URL

## v0.27.1
- CDR: Bugfixes

## v0.27.0
- [NEW] Initial Clinical Data Repository (CDR) support

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

