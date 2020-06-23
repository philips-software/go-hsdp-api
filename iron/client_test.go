package iron_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	muxIRON    *http.ServeMux
	serverIRON *httptest.Server
)

func setup(t *testing.T) func() {
	muxIRON = http.NewServeMux()
	serverIRON = httptest.NewServer(muxIRON)

	projectID := "48a0183d-a588-41c2-9979-737d15e9e860"

	muxIRON.HandleFunc("/projects/"+projectID+"/tasks", func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "GET", r.Method) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
  "tasks": [
    {
      "id": "bFp7OMpXdVsvRHp4sVtqb3gV",
      "created_at": "2020-06-23T09:47:07.967Z",
      "updated_at": "2020-06-23T10:19:58.119Z",
      "project_id": "Bny3gFLzLlMrFFDrujopyocu",
      "code_id": "5e6640a5fbce220009c0385e",
      "code_history_id": "5e6640a5fbce220009c0385f",
      "status": "cancelled",
      "msg": "Cancelled via API.",
      "code_name": "loafoe/siderite",
      "code_rev": "1",
      "start_time": "2020-06-23T09:47:11.85Z",
      "end_time": "0001-01-01T00:00:00Z",
      "timeout": 3600,
      "payload": "mu4xSCwztB79NcmrJvFEdRnw0priIxMDxLPencrypted",
      "schedule_id": "5eebb5113de052000a93b1f5",
      "message_id": "6841477577898197071",
      "cluster": "9PbpheKmd0bSHIelR7O6ChcH"
    },
    {
      "id": "ozAmEFk7mqs0UQXasmGQv2Js",
      "created_at": "2020-06-23T08:40:22.418Z",
      "updated_at": "2020-06-23T09:40:24.583Z",
      "project_id": "Bny3gFLzLlMrFFDrujopyocu",
      "code_id": "5e6640a5fbce220009c0385e",
      "code_history_id": "5e6640a5fbce220009c0385f",
      "status": "timeout",
      "msg": "Job timed out.",
      "code_name": "loafoe/siderite",
      "code_rev": "1",
      "start_time": "2020-06-23T08:40:24.44Z",
      "end_time": "2020-06-23T09:40:24.565Z",
      "duration": 3600125,
      "timeout": 3600,
      "payload": "mu4xSCwztB79NcmrJvFEdRnw0priIxMDxLPencrypted",
      "schedule_id": "5eebb5113de052000a93b1f5",
      "log_size": 34,
      "message_id": "6841460372259136626",
      "cluster": "9PbpheKmd0bSHIelR7O6ChcH"
    }`)
	})

	return func() {
		serverIRON.Close()
	}
}
