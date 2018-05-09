package logging

type Bundle struct {
	ResourceType string    `json:"resourceType"`
	Type         string    `json:"type"`
	Total        int       `json:"total"`
	ProductKey   string    `json:"productKey,omitempty"`
	Entry        []Element `json:"entry"`
}

type Element struct {
	Resource Resource `json:"resource"`
}

type Resource struct {
	ResourceType        string  `json:"resourceType"`        // LogEvent
	ID                  string  `json:"id"`                  // 7f4c85a8-e472-479f-b772-2916353d02a4
	ApplicationName     string  `json:"applicationName"`     // OPS
	EventID             string  `json:"eventId"`             // 110114
	Category            string  `json:"category"`            // TRACELOG
	Component           string  `json:"component"`           // "TEST"
	TransactionID       string  `json:"transactionId"`       // 2abd7355-cbdd-43e1-b32a-43ec19cd98f0
	ServiceName         string  `json:"serviceName"`         // OPS
	ApplicationInstance string  `json:"applicationInstance"` // INST-00002
	ApplicationVersion  string  `json:"applicationVersion"`  // 1.0.0
	OriginatingUser     string  `json:"originatingUser"`     // SomeUser
	ServerName          string  `json:"serverName"`          // app.example.com
	LogTime             string  `json:"logTime"`             // 2017-01-31T08:00:00Z
	Severity            string  `json:"severity"`            // INFO
	LogData             LogData `json:"logData"`
}

func (r *Resource) Valid() bool {
	if r.EventID != "" && r.TransactionID != "" && r.LogTime != "" && r.LogData.Message != "" {
		return true
	}
	return false
}

type LogData struct {
	Message string `json:"message"`
}
