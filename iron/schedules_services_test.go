package iron_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/iron"

	"github.com/stretchr/testify/assert"
)

func TestSchedulesServices_CreateSchedule(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	scheduleID := "bFp7OMpXdVsvRHp4sVtqb3gV"

	muxIRON.HandleFunc(client.Path("projects", projectID, "schedules"), func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "POST", r.Method) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"schedules":[{"id":"`+scheduleID+`"}]}`)
	})
	muxIRON.HandleFunc(client.Path("projects", projectID, "schedules", scheduleID), func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "GET", r.Method) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
      "id": "`+scheduleID+`",
      "created_at": "2020-06-25T11:32:46.72Z",
      "updated_at": "2020-06-25T11:32:46.72Z",
      "project_id": "`+projectID+`",
      "status": "scheduled",
      "code_name": "testandy",
      "start_at": "2020-06-25T11:32:46.72Z",
      "end_at": "0001-01-01T00:00:00Z",
      "next_start": "2020-06-25T11:32:46.72Z",
      "last_run_time": "0001-01-01T00:00:00Z",
      "timeout": 7200,
      "run_times": 3,
      "run_every": 3600,
      "cluster": "XKaaLazEd1sAUAyZZN8IG6Tg",
      "payload": "{}"
    }`)
	})

	schedule, resp, err := client.Schedules.CreateSchedule(iron.Schedule{
		CodeName: "foo",
		Payload:  "ron",
		Timeout:  7200,
		Cluster:  "XKaaLazEd1sAUAyZZN8IG6Tg",
	})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !assert.NotNil(t, schedule) {
		return
	}
	assert.Equal(t, scheduleID, schedule.ID)
}

func TestSchedulesServices_GetSchedule(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	scheduleID := "8GSI27QGIZ5sSYRiMBIoASz8"
	muxIRON.HandleFunc(client.Path("projects", projectID, "schedules", scheduleID), func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "GET", r.Method) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
      "id": "`+scheduleID+`",
      "created_at": "2020-06-25T11:32:46.72Z",
      "updated_at": "2020-06-25T11:32:46.72Z",
      "project_id": "`+projectID+`",
      "status": "scheduled",
      "code_name": "testandy",
      "start_at": "2020-06-25T11:32:46.72Z",
      "end_at": "0001-01-01T00:00:00Z",
      "next_start": "2020-06-25T11:32:46.72Z",
      "last_run_time": "0001-01-01T00:00:00Z",
      "timeout": 7200,
      "run_times": 3,
      "run_every": 3600,
      "cluster": "DRxYM4SCFZBiJrsWytWju38C",
      "payload": "{}"
    }`)
	})
	schedule, resp, err := client.Schedules.GetSchedule(scheduleID)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, schedule) {
		return
	}
	assert.Equal(t, scheduleID, schedule.ID)
}

func TestSchedulesServices_GetSchedules(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	scheduleID := "8GSI27QGIZ5sSYRiMBIoASz8"
	muxIRON.HandleFunc(client.Path("projects", projectID, "schedules"), func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "GET", r.Method) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
  "schedules": [
    {
      "id": "`+scheduleID+`",
      "created_at": "2020-06-25T11:32:46.72Z",
      "updated_at": "2020-06-25T11:32:46.72Z",
      "project_id": "`+projectID+`",
      "status": "scheduled",
      "code_name": "testandy",
      "start_at": "2020-06-25T11:32:46.72Z",
      "end_at": "0001-01-01T00:00:00Z",
      "next_start": "2020-06-25T11:32:46.72Z",
      "last_run_time": "0001-01-01T00:00:00Z",
      "timeout": 7200,
      "run_times": 3,
      "run_every": 3600,
      "cluster": "DRxYM4SCFZBiJrsWytWju38C",
      "payload": "{}"
    },
    {
      "id": "C8OvMrpP2f226nIMJQV5VNZz",
      "created_at": "2020-06-25T11:31:21.167Z",
      "updated_at": "2020-06-25T11:31:27.379Z",
      "project_id": "`+projectID+`",
      "status": "scheduled",
      "code_name": "testandy",
      "start_at": "2020-06-25T11:31:21.167Z",
      "end_at": "0001-01-01T00:00:00Z",
      "next_start": "2020-06-25T12:31:21.167Z",
      "last_run_time": "2020-06-25T11:31:27.334Z",
      "timeout": 7200,
      "run_times": 3,
      "run_count": 1,
      "run_every": 3600,
      "cluster": "DRxYM4SCFZBiJrsWytWju38C",
      "payload": "{}"
    }
  ]
}`)
	})
	tasks, resp, err := client.Schedules.GetSchedules()
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, tasks) {
		return
	}
	assert.Equal(t, 2, len(*tasks))
}

func TestSchedulesServices_CancelSchedule(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	scheduleID := "bFp7OMpXdVsvRHp4sVtqb3gV"

	muxIRON.HandleFunc(client.Path("projects", projectID, "schedules", scheduleID, "cancel"), func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "POST", r.Method) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"msg":"Cancelled"}`)
	})
	ok, resp, err := client.Schedules.CancelSchedule(scheduleID)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.True(t, ok) {
		return
	}
}
