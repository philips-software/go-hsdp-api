package tpns

import "github.com/philips-software/go-hsdp-api/fhir"

// MessagesService provides operations on TPNS messages
type MessagesService struct {
	client *Client
}

// Message describes a push message
type Message struct {
	MessageType      string            `json:"MessageType"`
	PropositionID    string            `json:"PropositionId"`
	CustomProperties map[string]string `json:"CustomProperties"`
	Lookup           bool              `json:"Lookup"`
	Content          string            `json:"Content"`
	Targets          []string          `json:"Targets"`
}

// Push pushes a message to a mobile client
func (m *MessagesService) Push(msg *Message) (bool, *Response, error) {
	req, err := m.client.NewTPNSRequest("POST", "tpns/PushMessage", msg, nil)
	if err != nil {
		return false, nil, err
	}

	var responseStruct fhir.IssueResponse

	resp, err := m.client.Do(req, &responseStruct)
	if err != nil {
		return false, resp, err
	}

	if err != nil {
		return false, nil, err
	}
	return true, resp, nil
}
