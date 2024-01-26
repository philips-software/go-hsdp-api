package dbs_test

import (
	"fmt"
	"github.com/philips-software/go-hsdp-api/connect/dbs"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func sqsBody(id, queueType string, infix string, status string) string {
	return fmt.Sprintf(`{
            "id": "%s",
            "meta": {
                "lastUpdated": "2022-12-06T10:18:11.947Z",
                "versionId": "1"
            },
            "name": "dbs-%s-%s",
            "description": "MySQSQueue",
            "status": "%s",
            "resourceType": "SQSSubscriber",
            "queueName": "dbs-%s-%s",
            "queueType": "%s",
            "deliveryDelaySeconds": 0,
            "messageRetentionPeriod": 345600,
            "receiveMessageWaitTimeSeconds": 0,
            "serverSideEncryption": true
}`, id, infix, id, status, infix, id, queueType)
}

func TestSQSCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	sqsID := "9f80f9e0-5cb2-4ebd-8980-03f550cb453f"
	queueType := "FIFO"
	infix := "my_infix"
	muxDBS.HandleFunc("/client-test/connect/databroker/Subscriber/SQS", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.Header().Set("Etag", "1")
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, sqsBody(sqsID, queueType, infix, "Creating"))
		}
	})
	muxDBS.HandleFunc("/client-test/connect/databroker/Subscriber/SQS/"+sqsID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, sqsBody(sqsID, queueType, infix, "Active"))
		case "PUT":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, sqsBody(sqsID, queueType, infix, "Updating"))
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		}
	})

	created, resp, err := dbsClient.Subscribers.CreateSQS(dbs.SQSSubscriberConfig{
		NameInfix:            infix,
		Description:          "MySQSQueue",
		QueueType:            queueType,
		ServerSideEncryption: true,
	})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, created) {
		return
	}
	assert.Equal(t, queueType, created.QueueType)
	assert.Equal(t, sqsID, created.ID)
	assert.NotNil(t, created.Status)

	res, resp, err := dbsClient.Subscribers.DeleteSQS(*created)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.True(t, res)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode())
}
