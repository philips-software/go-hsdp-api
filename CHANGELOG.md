# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## v0.63.0

- Update to Go 1.18

## v0.62.0

- IAM Service: support service updates (IAM March 2022 release)

## v0.61.5

- Connect MDM: workaround for query issue

## v0.61.4

- Logging: improve error message further

## v0.61.3

- Logging: also ammend server side resource rejection errors

## v0.61.2

- Logging: detailed error message for invalid resources 

# v0.61.1

- IAM: Fix bad typo

## v0.61.0

- IAM: don't ignore token errors. API breaking chnage!

## v0.60.1

- Cartel: fix add/remove user group calls

## v0.60.0

- Config: add apac2 (Tokyo) region

## v0.52.1

- DICOM: Fix response mismatch

## v0.52.0

- DICOM: Add notification support

## v0.51.8

- MDM: Improve create calls
- Log: mask JSON password fields
- Log: add UTC timestamp to debug log headers

## v0.51.3

- Core: print timestamps in debug log
- Core: mask Iron credentials payload

## v0.51.2

- Docker Service Keys: add fallback until fixed Docker API is deployed to all regions
- MDM: Fix application and proposition update calls

## v0.51.1

- Docker: Add GetLatestTag

## v0.51.0

- IAM: Some function signature changes for better error propagation

## v0.50.2

- Core: boolean validation fixes
- MDM: consistency tweaks

## v0.50.1

- MDM related bug fixes 

## v0.50.0

- NEW: Full Connect MDM support
- NEW: Host Service Discovery support

## v0.49.2

- Bump TDR API version to 5

## v0.49.1

- BREAKING change Logging: support true passthrough of LogEvent (d9d9014843f)

## v0.49.0

- Logging: refactor to use common debug logging code
- IAM Service: fix runaway recursive issue when refreshing

## v0.48.2

- Docker Registries: add additional fields to Registry struct

## v0.48.1

- Docker Namespaces: add update user access support

## v0.48.0

- Docker Registry: service keys, namespace and repository management

## v0.47.0

- Logging: Add traceId and spanId fields
- Console: add alerts structure
- IAM Service: recover from panicky pem.Decode
- CDR: Update delete extension URI
- IAM: Add new fields
- IAM: Ipdate LegacyUserUpdate to (undocumented) API v2
- IAM: update GetUser to v3 API
- IAM: support for preferredCommunicationChannel field

## v0.46.2

- CDR: R4 fixes

## v0.46.1

- CDL: Get all study pages
- CDL: Add GetStudyByTitle method
- IAM SMS Templates: locale and template fixes
- IAM SMS: Add If-Match pre-condition

## v0.46.0
 
- DICOM: Breaking API changes
- IAM SMS Gateway: Initial support
- IAM SMS Templates: Initial support

## v0.45.0
- 
- CDR: Support for R4

## v0.44.0

- DICOM Gateway fixes

## v0.43.0
 
- AI Inference: initial support
- AI Training: initial support
- AI Workspace: initial support

## v0.42.3
 
- CDL: add Export Routes CRD

### v0.42.2
 
- CDL: delete call for label definitions

### v0.42.1
 
- CDL: label definitions

### v0.42.0
 
- Improve error messages by returning the request body as part of the error

### v0.41.2

- CDL: add DTD support

### v0.41.1

- Security: migrate to github.com/golang-jwt/jwt

### v0.41.0

- Preliminary support for Clinical Data Lake (CDL) Resource studies

### v0.40.1

- Debug log filtering improvements
- Dependency upgrades
- DICOM Remote nodes API call fix

### v0.40.0

- Add Canada (ca1) region to service discovery
- Add vault-proxy service

### v0.39.1
- Fix notification related issues
- Improve masking of sensitive values in debug logs

### v0.39.0
- HSDP Notification support

### v0.38.1
- Filter known sensitive fields form debug logs
- Move cartel client to internal logger
- Move iron client to internal logger
- Fix DecryptPayload

### v0.38.0
- Add DecryptPayload to decrypt Iron payloads
- Export some convenience functions

### v0.37.0
- NEW: support for updating IAM service certificates
- Prepare for CDR changes
- Improve debug log output

### v0.36.6
- Add cn1 endpoints for IAM/IDM

### v0.36.5
- Add additional S3Creds regions

### v0.36.4
- Add all other possible actions

### v0.36.3
- Add ALL_BUCKET to S3 policy action list

### v0.36.2
- Dependency upgrade

### v0.36.1
- Tweak JSON structs

### v0.36.0
- NEW: Support immediate IAM account activation with optional password input
- Support user updates

### v0.35.6
- Wrap errors for better contextual errors

### v0.35.5
- Improved error checking for application and proposition creation

### v0.35.4
- Remove duplicate DICOM logging

### v0.35.3
- Expose TokenRefresh() in DICOM client

### v0.35.2
- Fix STL cert update call issue

### v0.35.1
- Add iam.Applications.GetApplicationByName()

### v0.35.0
- NEW: Secure Transport Layer (STL) support

### v0.34.4
- Fix some PKI methods
- Add IAM token revoke calls

### v0.34.3
- Remove elastic due to license change

### v0.34.2
- Add pki.Services.GetCertificates() method

### v0.34.1
- Fix minimum version in go.mod

### v0.34.0
- Remove online config refresh code
- Use Go 1.16 embed to bundle hsdp.json
- This version only works with Go 1.16+ (breaking change!)

### v0.33.0
- NEW: DICOM Config API support
- Create internal package for consistent better versioning
- S3Creds renaming (breaking change!)

### v0.32.3
- Disable keep-alive for Cartel

### v0.32.2
- Better error reporting in Cartel

### v0.32.1
- Add cartel.BastionHost() 

### v0.32.0
- Proxy support

### v0.31.2
- Email template: change locale handling

### v0.31.1
- Bug fixes

### v0.31.0
- [NEW] IDM Email Templates support

### v0.30.0
- Remove fhir package

### v0.29.1
- Fix HSDP Audit documentation
- Unexport many methods

### v0.29.0
- [NEW] HSDP Audit support

### v0.28.0
- Logging: support bearer token auth (client credentials / service identities)
- Tweaking of CDR client

### v0.27.2
- CDR: Add endpoint URL

### v0.27.1
- CDR: Bugfixes

### v0.27.0
- [NEW] Initial Clinical Data Repository (CDR) support

### v0.26.0
- [Breaking] Refactored config API. Master data is stored as JSON
- Remove toml dependency

### v0.25.0
- [NEW] HSDP PKI API support

### v0.24.0
- Logging fixes
- [Cartel] Minor API tweaks for better error reporting

### v0.23.0
- Maintenance release
- Add Tag() option for Cartel

### v0.22.0
- [NEW] Console API support: Autoscalers

### v0.21.1
- Fix config URL

### v0.21.0
- [NEW] Hosted appstream (HAS) support
- [NEW] IAM Password policies management

### v0.20.0
- Autoconfig support for IAM, Cartel and Logging
- Add ap3 region
- Fallback mechanism for discovery

### v0.19.0
- [NEW] Service/config discovery
- Share structs with gautocloud-connectors

### v0.18.0
- [NEW] IronIO Worker support
  Codes CRUD
  Tasks CRUD
  Schedules CRUD
  Clusters Read/stats

### v0.17.0
- [IAM] Metadata field changes

### v0.16.0
- [IAM] Add UpdateClient() method 
- [Logging] Add Meta tags to Resource

### v0.15.0

- [NEW] Cartel API support
- [IAM] Switch to SCIM based organization management
- [IAM] Add Organizations.DeleteOrganization()
- [IAM] Add Organizations.DeleteStatus()
- [Logging] Detect errors in batch sends

### v0.14.0
- [IAM] Move user find API to v2
- [IAM] Update go-hsdp-signer
- [Logging] Better support for custom logging

### v0.2.0
- Upgrade github.com/Jeffail/gabs

### v0.1.0
- IAM support
- Logging support
- TDR basic support

[Unreleased]: https://github.com/philips-software/go-hsdp-api/compare/1.0.0...HEAD

- MDM: Support for bootstrap OAuth client scopes
