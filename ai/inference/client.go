package inference

import (
	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/ai"
	"github.com/philips-software/go-hsdp-api/iam"
)

// A Client manages communication with HSDP AI-Inference API
type Client struct {
	*ai.Client
	ComputeEnvironment *ai.ComputeEnvironmentService
	Model              *ModelService
	Job                *ai.JobService
}

// NewClient returns a new HSDP AI-Inference API client. A configured IAM client
// must be provided as the underlying API requires an IAM token
func NewClient(iamClient *iam.Client, config *ai.Config) (*Client, error) {
	config.Service = "inference"
	client, err := ai.NewClient(iamClient, config)
	if err != nil {
		return nil, err
	}
	inferenceClient := &Client{
		Client:             client,
		Model:              &ModelService{client: client, validate: validator.New()},
		Job:                &ai.JobService{Client: client, Validate: validator.New(), Path: "JobService"},
		ComputeEnvironment: &ai.ComputeEnvironmentService{Client: client, Validate: validator.New()},
	}

	return inferenceClient, nil
}
