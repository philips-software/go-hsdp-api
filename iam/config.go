package iam

import (
	hsdpsigner "github.com/philips-software/go-hsdp-signer"
)

// Config contains the configuration of a client
type Config struct {
	Region           string
	Environment      string
	OAuth2ClientID   string
	OAuth2Secret     string
	SharedKey        string
	SecretKey        string
	BaseIAMURL       string
	BaseIDMURL       string
	OrgAdminUsername string
	OrgAdminPassword string
	IAMURL           string
	IDMURL           string
	Scopes           []string
	RootOrgID        string
	Debug            bool
	DebugLog         string
	Signer           *hsdpsigner.Signer
}
