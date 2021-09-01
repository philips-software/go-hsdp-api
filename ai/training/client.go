package training

import (
	"github.com/go-playground/validator/v10"
	"github.com/philips-software/go-hsdp-api/ai"
	"github.com/philips-software/go-hsdp-api/iam"
)

// A Client manages communication with HSDP AI-Training API
type Client struct {
	*ai.Client
	ComputeEnvironment *ai.ComputeEnvironmentService
	Job                *ai.JobService
}

// NewClient returns a new HSDP AI-Training API client. A configured IAM client
// must be provided as the underlying API requires an IAM token
func NewClient(iamClient *iam.Client, config *ai.Config) (*Client, error) {
	config.Service = "training"
	client, err := ai.NewClient(iamClient, config)
	if err != nil {
		return nil, err
	}
	trainingClient := &Client{
		Client:             client,
		Job:                &ai.JobService{Client: client, Validate: validator.New(), Path: "Job"},
		ComputeEnvironment: &ai.ComputeEnvironmentService{Client: client, Validate: validator.New()},
	}
	return trainingClient, nil
}
