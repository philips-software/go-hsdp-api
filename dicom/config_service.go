package dicom

import (
	"github.com/google/fhir/go/jsonformat"
)

// ConfigService
type ConfigService struct {
	client  *Client
	profile string
	ma      *jsonformat.Marshaller
	um      *jsonformat.Unmarshaller
}

// QueryOptions holds optional query options for requests
type QueryOptions struct {
	OrganizationID *string `url:"organizationId,omitempty"`
}
