package notification_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/notification"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptionCRD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	storeID := "f65b7642-442e-4597-a64f-260f9251ca1d"

	someSubscription := notification.Subscription{
		TopicID:              "some-topic-id",
		SubscriberID:         "some-subscriber-id",
		SubscriptionEndpoint: "https://foo.bar/endpoint",
	}

	muxNotification.HandleFunc("/core/notification/Subscription", func(w http.ResponseWriter, r *http.Request) {
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
			var received notification.Subscription
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
			_, _ = io.WriteString(w, `{
  "resourceType": "Bundle",
  "type": "searchset",
  "total": 1,
  "link": [
    {
      "relation": "self",
      "url": "https://notification-dev.us-east.philips-healthsuite.com/core/notification/Subscription"
    }
  ],
  "entry": [
    {
      "topic_id": "some-topic-id",
      "subscriber_id": "some-subscriber-id",
      "subscription_endpoint": "https://foo.bar/endpoint",
      "_id": "`+storeID+`",
      "resourceType": "Subscription"
    }
  ]
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	muxNotification.HandleFunc("/core/notification/Subscription/"+storeID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	muxNotification.HandleFunc("/core/notification/Subscription/_confirm", func(w http.ResponseWriter, r *http.Request) {
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
			var received notification.ConfirmRequest
			err := json.NewDecoder(r.Body).Decode(&received)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var response = notification.Subscription{
				ID:                   storeID,
				SubscriptionARN:      received.TopicARN,
				TopicID:              "random-id",
				SubscriberID:         "random-id",
				ResourceType:         "Subscription",
				SubscriptionEndpoint: "https://notification-receiver.bogus/notification/Test-1-27",
			}
			resp, err := json.Marshal(&response)
			if !assert.Nil(t, err) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, string(resp))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	created, resp, err := notificationClient.Subscription.CreateSubscription(someSubscription)
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
	ok, resp, err := notificationClient.Subscription.DeleteSubscription(notification.Subscription{ID: storeID})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.True(t, ok)

	item, resp, err := notificationClient.Subscription.GetSubscription(storeID)
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

	item, resp, err = notificationClient.Subscription.ConfirmSubscription(notification.ConfirmRequest{
		Token:    "some-token",
		TopicARN: "arn:thing",
		Endpoint: "https://some.bogus/endpoint",
	})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, item) {
		return
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, storeID, item.ID)
}
