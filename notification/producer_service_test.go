package notification_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/notification"
	"github.com/stretchr/testify/assert"
)

func TestProducersCRD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	storeID := "f65b7642-442e-4597-a64f-260f9251ca1d"
	orgID := "614b0053-7a57-44d8-ba8a-809b9362a9a6"

	someProducer := notification.Producer{
		ManagingOrganizationID:      orgID,
		ManagingOrganization:        "SomeOrg",
		ProducerProductName:         "test",
		ProducerServiceName:         "test",
		ProducerServiceInstanceName: "test",
		ProducerServiceBaseURL:      "https://foo",
		ProducerServicePathURL:      "/bar",
	}

	muxNotification.HandleFunc("/core/notification/Producer", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			if !assert.Equal(t, "application/json", r.Header.Get("Content-Type")) {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				return
			}
			if !assert.Equal(t, notification.APIVersion, r.Header.Get("API-Version")) {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
			var receivedProducer notification.Producer
			err := json.NewDecoder(r.Body).Decode(&receivedProducer)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			receivedProducer.ID = storeID
			resp, err := json.Marshal(&receivedProducer)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, _ = io.WriteString(w, string(resp))
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "resourceType": "Bundle",
  "type": "searchset",
  "total": 1,
  "link": [
    {
      "relation": "self",
      "url": "https://notification-dev.us-east.philips-healthsuite.com/core/notification/Producer"
    }
  ],
  "entry": [
    {
      "managingOrganization": "DssSmokeTest",
      "producerProductName": "exampleProduct",
      "producerServiceName": "exampleServiceName",
      "producerServiceInstanceName": "exampleServiceInstance",
      "producerServiceBaseUrl": "https://ns-producer.cloud.pcftest.com/",
      "producerServicePathUrl": "notification/create",
      "description": "product description",
      "managingOrganizationId": "`+orgID+`",
      "_id": "`+storeID+`",
      "resourceType": "Producer"
    }
  ]
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	muxNotification.HandleFunc("/core/notification/Producer/"+storeID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	created, resp, err := notificationClient.Producer.CreateProducer(someProducer)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, created) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, created.ID, storeID)
	ok, resp, err := notificationClient.Producer.DeleteProducer(notification.Producer{ID: storeID})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.True(t, ok)

	item, resp, err := notificationClient.Producer.GetProducer(storeID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, item) {
		return
	}
	assert.Equal(t, storeID, item.ID)
}
