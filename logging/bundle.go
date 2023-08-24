package logging

import (
	"encoding/json"
	"fmt"
)

// Bundle is a FHIR bundle resource
// There is just enough there to create the logging payload
type Bundle struct {
	ResourceType string    `json:"resourceType"`
	Type         string    `json:"type"`
	Total        int       `json:"total"`
	ProductKey   string    `json:"productKey,omitempty"`
	Entry        []Element `json:"entry"`
}

// Element is a FHIR element resource
type Element struct {
	Resource Resource `json:"resource"`
}

// Resource is a logging resource
type Resource struct {
	ResourceType        string                 `json:"resourceType"`          // LogEvent
	ID                  string                 `json:"id"`                    // 7f4c85a8-e472-479f-b772-2916353d02a4
	ApplicationName     string                 `json:"applicationName"`       // OPS
	EventID             string                 `json:"eventId"`               // 110114
	Category            string                 `json:"category"`              // TRACELOG
	Component           string                 `json:"component"`             // "TEST"
	TransactionID       string                 `json:"transactionId"`         // 2abd7355-cbdd-43e1-b32a-43ec19cd98f0
	ServiceName         string                 `json:"serviceName,omitempty"` // OPS
	ApplicationInstance string                 `json:"applicationInstance"`   // INST-00002
	ApplicationVersion  string                 `json:"applicationVersion"`    // 1.0.0
	OriginatingUser     string                 `json:"originatingUser"`       // SomeUser
	ServerName          string                 `json:"serverName"`            // app.example.com
	LogTime             string                 `json:"logTime"`               // 2017-01-31T08:00:00Z
	Severity            string                 `json:"severity"`              // INFO
	TraceID             string                 `json:"traceId,omitempty"`     // xxx
	SpanID              string                 `json:"spanId,omitempty"`      // yyy
	LogData             LogData                `json:"logData"`               // Log data
	Custom              json.RawMessage        `json:"custom,omitempty"`      // Custom log fields
	Meta                map[string]interface{} `json:"-"`
	Error               error                  `json:"-"`
}

// LogData is the payload of a log message
type LogData struct {
	Message string `json:"message"`
}

// Valid returns true if a resource is valid according to HSDP rules, false otherwise
func (r *Resource) Valid() bool {
	var u map[string]interface{}

	if r.EventID == "" {
		r.Error = fmt.Errorf("EventID field is blank")
		return false
	}

	if r.TransactionID == "" {
		r.Error = fmt.Errorf("TransactionID field is blank")
		return false
	}
	if r.LogTime == "" {
		r.Error = fmt.Errorf("LogTime field is blank")
		return false
	}
	if r.LogData.Message == "" {
		r.Error = fmt.Errorf("LogData.Message field is blank")
		return false
	}
	if len(r.Custom) > 0 {
		if err := json.Unmarshal(r.Custom, &u); err != nil {
			r.Error = fmt.Errorf("custom field unmarshal error: %w", err)
			return false
		}
	}
	return true
}

type bundleErrorResponse struct {
	Issue []struct {
		Severity string `json:"severity"`
		Code     string `json:"code"`
		Details  struct {
			Coding []struct {
				System string `json:"system"`
				Code   string `json:"code"`
			} `json:"coding"`
			Text string `json:"text"`
		} `json:"details"`
		Diagnostics string   `json:"diagnostics"`
		Location    []string `json:"location"`
	} `json:"issue"`
	ResourceType string `json:"resourceType"`
}
