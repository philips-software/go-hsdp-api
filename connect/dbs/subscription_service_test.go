package dbs_test

import (
	"fmt"
	"github.com/philips-software/go-hsdp-api/connect/dbs"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func subscriptionBody(id, dataType string, infix string, status string, subscriberId string) string {
	return fmt.Sprintf(`{
			"dataType": "%s",
            "description": "TestSubscription",
            "deliverDataOnly": false,
            "id": "%s",
            "meta": {
                "lastUpdated": "2022-12-06T10:23:20.819Z",
                "versionId": "1"
            },
            "name": "dbs_d77e1fb8-9e13-48e0-9d20-b6365793a181_%s_%s",
            "resourceType": "TopicSubscription",
            "ruleName": "dbs_d77e1fb8-9e13-48e0-9d20-b6365793a181_%s_%s",
            "status": "%s",
            "subscriber": {
                "id": "%s",
                "type": "SQSSubscriber",
                "location": "https://databroker-client-test.eu01.connect.hsdp.io/client-test/connect/databroker/Subscriber/SQS/%s"
            }
}`, dataType, id, infix, id, infix, id, status, subscriberId, subscriberId)
}

func TestSubscriptionCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	sqsID := "9f80f9e0-5cb2-4ebd-8980-03f550cb453f"
	subscriptionID := "1ca7251b-42a1-4560-99a5-2b359c6f3914"
	infix := "my_infix"
	dataType := "my-datatype"
	muxDBS.HandleFunc("/client-test/connect/databroker/Subscription/Topic", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.Header().Set("Etag", "1")
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, subscriptionBody(subscriptionID, dataType, infix, "Creating", sqsID))
		}
	})
	muxDBS.HandleFunc("/client-test/connect/databroker/Subscription/Topic/"+subscriptionID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, subscriptionBody(subscriptionID, dataType, infix, "Active", sqsID))
		case "PUT":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, subscriptionBody(subscriptionID, dataType, infix, "Updating", sqsID))
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		}
	})

	created, resp, err := dbsClient.Subscriptions.CreateTopicSubscription(dbs.TopicSubscriptionConfig{
		NameInfix:    infix,
		Description:  "MyTopicSubscription",
		SubscriberId: sqsID,
		DataType:     dataType,
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
	assert.Equal(t, dataType, created.DataType)
	assert.Equal(t, sqsID, created.Subscriber.ID)
	assert.Equal(t, subscriptionID, created.ID)
	assert.NotNil(t, created.Status)

	res, resp, err := dbsClient.Subscriptions.DeleteTopicSubscription(*created)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.True(t, res)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode())
}
