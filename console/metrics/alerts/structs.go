package alerts

import (
	"time"
)

// Payload is the JSON format of HSDP Metrics webhook alerts
type Payload struct {
	Receiver          string     `json:"receiver"`
	Status            string     `json:"status"`
	Alerts            []Alert    `json:"alerts"`
	GroupLabels       GroupLabel `json:"groupLabels"`
	CommonLabels      Label      `json:"commonLabels"`
	CommonAnnotations Annotation `json:"commonAnnotations"`
	ExternalURL       string     `json:"externalURL"`
	Version           string     `json:"version"`
	GroupKey          string     `json:"groupKey"`
}

// Alert describes an alert
type Alert struct {
	Status       string     `json:"status"`
	Labels       Label      `json:"labels"`
	Annotations  Annotation `json:"annotations"`
	StartsAt     time.Time  `json:"startsAt"`
	EndsAt       time.Time  `json:"endsAt"`
	GeneratorURL string     `json:"generatorURL"`
}

// Annotation describes an alert
type Annotation struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

// Label describes a label
type Label struct {
	AlertName      string `json:"alertname"`
	Application    string `json:"application"`
	ApplicationID  string `json:"application_id"`
	BrokerID       string `json:"broker_id"`
	Deployment     string `json:"deployment"`
	Instance       string `json:"instance"`
	Job            string `json:"job"`
	Organization   string `json:"organization"`
	OrganizationID string `json:"organization_id"`
	Region         string `json:"region"`
	Severity       string `json:"severity"`
	Space          string `json:"space"`
	SpaceID        string `json:"space_id"`
}

// GroupLabel describes a group label
type GroupLabel struct {
	AlertName   string `json:"alertname"`
	Application string `json:"application"`
}
