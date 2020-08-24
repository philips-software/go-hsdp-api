package has

// The HAS Resource API provides access to a windows application virtualization service.
// This API allows a user to manage a pool of storage and computing resources used for session based access.
// The pool is managed by a scheduler, which will take care of day-to-day resource scheduling, and which can
// be influenced or overruled through parameters in this interface.
//
// Ref: https://www.hsdp.io/documentation/hosted-application-streaming/

import (
	"net/http"
)

// ContractsService provides operations on TDR contracts
type ResourcesService struct {
	orgID  string
	client *Client
}

// Constants
const (
	HASAPIVersion = "1"
)

// Resource
type Resource struct {
	// the id of the resource.
	ID string `json:"id,omitempty"`
	// id for virtual machine image. This typically lives as an AMI and is defined elsewhere. Example: "ami-454654654"
	ImageID string `json:"imageId" validate:"required"`
	// Reference to resource type of cloud provider. Default will be a machine type configured for a VM image. Example: "g3.4xlarge"
	ResourceType string `json:"resourceType" validate:"required"`
	// The region where the resource will be hosted. See the providers documentation for available regions.
	Region string `json:"region" validate:"required"`
	// Number of resources requested. The maximum amount to request resources in one api call is 10. Defaults to 1.
	Count int `json:"count" validate:"min=1,max=10,required"`
	// Allows to perform filters and queries on a cluster of resources.
	ClusterTag string `json:"clusterTag" validate:"required"`
	EBS        EBS    `json:"ebs" validate:"required"`
	// Unique identifier to a resource (i.e. server/storage instance)
	ResourceID string `json:"resourceId,omitempty"`
	// The id of the organization that this resource belongs to.
	OrganizationID string `json:"organizationId,omitempty"`
	// The current sessionId of the session which has claimed the resource, will be not present when not claimed.
	SessionID string `json:"sessionId,omitempty"`
	// State of AWS resource.
	State string `json:"state,omitempty" enum:"PENDING|RUNNING|SHUTTING-DOWN|TERMINATED|REBOOTING|STOPPING|STOPPED|UNKNOWN"`
	// The DNS name of the resource
	DNS string `json:"dns,omitempty"`
	// When a resource is disabled it means that the resource is deleted and wont be used for sessions.
	// It will only be used as reference for historical session.
	Disabled bool `json:"disabled,omitempty"`
}

// EBS options provided to AWS
type EBS struct {
	DeleteOnTermination bool   `json:"DeleteOnTermination"`
	Encrypted           bool   `json:"Encrypted"`
	Iops                int    `json:"Iops,omitempty"`
	KmsKeyID            string `json:"KmsKeyId,omitempty"`
	SnapshotID          string `json:"SnapshotId,omitempty"`
	VolumeSize          int    `json:"VolumeSize" validate:"min=20,required"`
	VolumeType          string `json:"VolumeType" validate:"required" enum:"standard|io1|gp2|sc1|st1"`
}

type resourceResponse struct {
	Resources        []Resource `json:"resources"`
	Error            string     `json:"error,omitempty"`
	ErrorDescription string     `json:"error_description,omitempty"`
}

// Result of a resourc action API call
type Result struct {
	ResourceID    string `json:"resourceId"`
	Action        string `json:"action"`
	ResultCode    int    `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

// ResourcesReports represents a list of Result values
type ResourcesReport struct {
	Results []Result `json:"results"`
}

// CreateResource creates a new resource in HAS
// A server will be prepared immediately.
// That is, an EBS backed AMI will be created from provided reference and started.
// The Post operation allows a user to add new resources to the pool or cluster.
// This is an operational action, and requires elevated permissions to operate.
// This endpoint requires HAS_RESOURCE.ALL permission.
func (c *ResourcesService) CreateResource(resource Resource) (*[]Resource, *Response, error) {
	req, err := c.client.NewHASRequest("POST", "resource", &resource, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("organizationId", c.orgID)
	req.Header.Set("Api-Version", HASAPIVersion)

	var cr resourceResponse
	resp, err := c.client.Do(req, &cr)
	if err != nil {
		return nil, resp, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp, err
	}
	if len(cr.Resources) == 0 {
		return nil, resp, ErrEmptyResult
	}
	return &cr.Resources, resp, nil
}

type ResourceOptions struct {
	ClusterTag   *string `url:"clusterTag,omitempty"`
	ImageID      *string `url:"imageId,omitempty"`
	ResourceType *string `url:"resourceType,omitempty"`
	ResourceID   *string `url:"resourceId,omitempty"`
	SessionID    *string `url:"sessionId,omitempty"`
	State        *string `url:"state,omitempty"`
	Region       *string `url:"region,omitempty"`
	Force        *bool   `url:"force,omitempty"`
}

// GetResources retrieves resources in HAS
// Get overview of all current resource claims.
// All fields are case-sensitive.
// This endpoint requires HAS_RESOURCE.ALL permission.
func (c *ResourcesService) GetResources(opt *ResourceOptions, options ...OptionFunc) (*[]Resource, *Response, error) {
	req, err := c.client.NewHASRequest("GET", "resource", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("organizationId", c.orgID)
	req.Header.Set("Api-Version", HASAPIVersion)

	var gr resourceResponse

	resp, err := c.client.Do(req, &gr)
	if err != nil {
		return nil, resp, err
	}
	return &gr.Resources, resp, nil
}

// DeleteResources deletes resources in HAS
// Delete multiple resources that are reserved.
// Resources that are claimed by a user can only
// be deleted by adding the force option.
// This endpoint requires HAS_RESOURCE.ALL permission.
func (c *ResourcesService) DeleteResources(opt *ResourceOptions, options ...OptionFunc) (*ResourcesReport, *Response, error) {
	req, err := c.client.NewHASRequest("DELETE", "resource", opt, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("organizationId", c.orgID)
	req.Header.Set("Api-Version", HASAPIVersion)

	var dr ResourcesReport

	resp, err := c.client.Do(req, &dr)
	if err != nil {
		return nil, resp, err
	}
	return &dr, resp, nil
}

// StartResource starts a resource in HAS
// Start a resource to make sure it is in Running state.
// This endpoint requires HAS_RESOURCE.ALL permission.
func (c *ResourcesService) StartResource(resources []string, options ...OptionFunc) (*[]Resource, *Response, error) {
	return c.startStopResource("start", resources, options...)
}

// StopResource stops a resource in HAS
// Stop a resource to make sure it is in Stopped state.
// This endpoint requires HAS_RESOURCE.ALL permission.
func (c *ResourcesService) StopResource(resources []string, options ...OptionFunc) (*[]Resource, *Response, error) {
	return c.startStopResource("stop", resources, options...)
}

func (c *ResourcesService) startStopResource(action string, resources []string, options ...OptionFunc) (*[]Resource, *Response, error) {
	var resourceList = struct {
		ResourceIDs []string `json:"resourceIds"`
	}{
		ResourceIDs: resources,
	}

	req, err := c.client.NewHASRequest("POST", "resource/"+action, &resourceList, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("organizationId", c.orgID)
	req.Header.Set("Api-Version", HASAPIVersion)

	var gr resourceResponse

	resp, err := c.client.Do(req, &gr)
	if err != nil {
		return nil, resp, err
	}
	return &gr.Resources, resp, nil
}

// GetResource retrieves a resource in HAS
// Get overview of the requested resource.
// This endpoint requires HAS_RESOURCE.ALL permission.
func (c *ResourcesService) GetResource(resourceID string, options ...OptionFunc) (*Resource, *Response, error) {
	req, err := c.client.NewHASRequest("GET", "resource/"+resourceID, nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("organizationId", c.orgID)
	req.Header.Set("Api-Version", HASAPIVersion)

	var gr Resource

	resp, err := c.client.Do(req, &gr)
	if err != nil {
		return nil, resp, err
	}
	return &gr, resp, nil
}

// DeleteResource deletes resources in HAS
// Delete a resource that is reserved.
// Resources that are claimed by a user can only be
// deleted by adding the force option.
// This endpoint requires HAS_RESOURCE.ALL permission.
func (c *ResourcesService) DeleteResource(resourceID string, options ...OptionFunc) (*ResourcesReport, *Response, error) {
	req, err := c.client.NewHASRequest("DELETE", "resource/"+resourceID, nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("organizationId", c.orgID)
	req.Header.Set("Api-Version", HASAPIVersion)

	var dr ResourcesReport

	resp, err := c.client.Do(req, &dr)
	if err != nil {
		return nil, resp, err
	}
	return &dr, resp, nil
}
