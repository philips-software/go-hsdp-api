package logging

type DHPLogMessage struct {
	Category            string     `json:"cat"`
	EventID             string     `json:"evt"`
	ApplicationVersion  string     `json:"ver"`
	Component           string     `json:"cmp"`
	ApplicationName     string     `json:"app"`
	ApplicationInstance string     `json:"inst"`
	ServerName          string     `json:"srv"`
	TransactionID       string     `json:"trns"`
	ServiceName         string     `json:"service"`
	LogTime             string     `json:"time"`
	OriginatingUser     string     `json:"usr"`
	Severity            string     `json:"sev"`
	LogData             DHPLogData `json:"val"`
}

type DHPLogData struct {
	Message string `json:"message"`
}

func (m *DHPLogMessage) Valid() bool {
	if m.LogData.Message != "" && m.LogTime != "" && m.TransactionID != "" {
		return true
	}
	return false

}
