package notification_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/notification"
	"github.com/stretchr/testify/assert"
)

func TestSubscribersCRD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	storeID := "f65b7642-442e-4597-a64f-260f9251ca1d"
	orgID := "614b0053-7a57-44d8-ba8a-809b9362a9a6"

	someSubscriber := notification.Subscriber{
		ManagingOrganizationID:        orgID,
		ManagingOrganization:          "SomeOrg",
		SubscriberProductName:         "test",
		SubscriberServicename:         "test",
		SubscriberServiceinstanceName: "test",
		SubscriberServiceBaseURL:      "https://foo",
		SubscriberServicePathURL:      "/bar",
	}

	muxNotification.HandleFunc("/core/notification/Subscriber", func(w http.ResponseWriter, r *http.Request) {
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
			var received notification.Subscriber
			err := json.NewDecoder(r.Body).Decode(&received)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			received.ID = storeID
			resp, err := json.Marshal(&received)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, _ = io.WriteString(w, string(resp))
		case "GET":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{
  "resourceType": "Bundle",
  "type": "searchset",
  "total": 1,
  "link": [
    {
      "relation": "self",
      "url": "https://notification-dev.us-east.philips-healthsuite.com/core/notification/Subscriber"
    }
  ],
  "entry": [
    {
      "managingOrganization": "DssSmokeTest",
      "subscriberProductName": "subsciberProd",
      "subscriberServiceName": "subsciberService",
      "subscriberServiceInstanceName": "serviceInsttest12",
      "subscriberServiceBaseUrl": "https://ns-client-logdev.cloud.pcftest.com/",
      "subscriberServicePathUrl": "core",
      "description": "subscriber description",
      "managingOrganizationId": "`+orgID+`",
      "_id": "`+storeID+`",
      "resourceType": "Subscriber"
    }
  ]
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	muxNotification.HandleFunc("/core/notification/Subscriber/"+storeID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	created, resp, err := notificationClient.Subscriber.CreateSubscriber(someSubscriber)
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
	ok, resp, err := notificationClient.Subscriber.DeleteSubscriber(notification.Subscriber{ID: storeID})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.True(t, ok)

	item, resp, err := notificationClient.Subscriber.GetSubscriber(storeID)
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
