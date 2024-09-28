package notification_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/notification"
	"github.com/stretchr/testify/assert"
)

func TestTopicCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	storeID := "f65b7642-442e-4597-a64f-260f9251ca1d"

	someTopic := notification.Topic{
		Name:          "name",
		Scope:         "public",
		ProducerID:    "some-producer-id",
		AllowedScopes: []string{"*.*.*.*"},
		IsAuditable:   false,
		Description:   "Some description",
	}

	muxNotification.HandleFunc("/core/notification/Topic", func(w http.ResponseWriter, r *http.Request) {
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
			var received notification.Topic
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
			_, _ = w.Write(resp)
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "resourceType": "Bundle",
  "type": "searchset",
  "total": 1,
  "link": [
    {
      "relation": "self",
      "url": "https://notification-dev.us-east.philips-healthsuite.com/core/notification/Topic"
    }
  ],
  "entry": [
    {
      "name": "Topic1",
      "producerId": "6eacf99c-704b-4c82-b966-1c770f2333c8",
      "scope": "public",
      "description": "topic description",
      "allowedScopes": [
        "*.*.*.NotificationTest"
      ],
      "isAuditable": true,
      "_id": "`+storeID+`",
      "resourceType": "Topic"
    }
  ]
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	muxNotification.HandleFunc("/core/notification/Topic/"+storeID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	created, resp, err := notificationClient.Topic.CreateTopic(someTopic)
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
	ok, resp, err := notificationClient.Topic.DeleteTopic(notification.Topic{ID: storeID})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.True(t, ok)

	item, resp, err := notificationClient.Topic.GetTopic(storeID)
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
