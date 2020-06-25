package iron_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/iron"
	"github.com/stretchr/testify/assert"
)

func TestCodesServices_CreateOrUpdateCode(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	codeID := "K6hyfuQzEmB9tDnKKHbKljjr"
	muxIRON.HandleFunc(client.Path("projects", projectID, "codes"), func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "POST", r.Method) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"id":"`+codeID+`"}`)
	})
	muxIRON.HandleFunc(client.Path("projects", projectID, "codes", codeID), func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "GET", r.Method) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
  "id": "`+codeID+`",
  "created_at": "2020-06-23T23:13:43.949Z",
  "project_id": "5e20da41d748ad000ace7654",
  "name": "testandy",
  "image": "loafoe/siderite:0.99.20",
  "rev": 2,
  "latest_history_id": "5ef3a3c96f3bb20009ba9952",
  "latest_change": "2020-06-24T19:04:41.782Z",
  "archived_at": "0001-01-01T00:00:00Z"
}`)
	})

	code, resp, err := client.Codes.CreateOrUpdateCode(iron.Code{
		Name:  "testandy",
		Image: "loafoe/siderite:0.99.20",
	})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !assert.NotNil(t, code) {
		return
	}
	assert.Equal(t, codeID, code.ID)
}

func TestCodesServices_GetCodes(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	codeID := "o30Zb9XLoDKYOfn721JqDTel"
	muxIRON.HandleFunc(client.Path("projects", projectID, "codes"), func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "GET", r.Method) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
  "codes": [
    {
      "id": "`+codeID+`",
      "created_at": "2020-03-17T13:57:47.383Z",
      "project_id": "5e20da41d748ad000ace7654",
      "name": "docker.na1.hsdp.io/loafoe/iron-streaming-backup",
      "stack": "",
      "image": "docker.na1.hsdp.io/loafoe/iron-streaming-backup:0.0.1",
      "runtime": "",
      "command": "",
      "rev": 7,
      "latest_history_id": "o30Zb9XLoDKYOfn721JqDTel",
      "latest_change": "2020-03-17T14:20:33.723Z",
      "archived_at": "0001-01-01T00:00:00Z"
    },
    {
      "id": "U04fyVshO2u1htvdBx8l6SYm",
      "created_at": "2020-03-27T11:51:52.059Z",
      "project_id": "5e20da41d748ad000ace7654",
      "name": "loafoe/iron-streaming-backup",
      "stack": "",
      "image": "loafoe/iron-streaming-backup:latest",
      "runtime": "",
      "command": "",
      "rev": 1,
      "latest_history_id": "U04fyVshO2u1htvdBx8l6SYm",
      "latest_change": "2020-03-27T11:51:52.059Z",
      "archived_at": "0001-01-01T00:00:00Z"
    }
]}`)
	})

	codes, resp, err := client.Codes.GetCodes()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, codes) {
		return
	}
	if !assert.Equal(t, 2, len(*codes)) {
		return
	}
	assert.Equal(t, codeID, (*codes)[0].ID)
}

func TestCodesServices_GetCode(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	codeID := "o30Zb9XLoDKYOfn721JqDTel"
	muxIRON.HandleFunc(client.Path("projects", projectID, "codes", codeID), func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "GET", r.Method) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
      "id": "`+codeID+`",
      "created_at": "2020-03-17T13:57:47.383Z",
      "project_id": "5e20da41d748ad000ace7654",
      "name": "docker.na1.hsdp.io/loafoe/iron-streaming-backup",
      "stack": "",
      "image": "docker.na1.hsdp.io/loafoe/iron-streaming-backup:0.0.1",
      "runtime": "",
      "command": "",
      "rev": 7,
      "latest_history_id": "o30Zb9XLoDKYOfn721JqDTel",
      "latest_change": "2020-03-17T14:20:33.723Z",
      "archived_at": "0001-01-01T00:00:00Z"
    }`)
	})

	code, resp, err := client.Codes.GetCode(codeID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, code) {
		return
	}
	assert.Equal(t, codeID, code.ID)
}

func TestCodesServices_DeleteCode(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	codeID := "bFp7OMpXdVsvRHp4sVtqb3gV"

	muxIRON.HandleFunc(client.Path("projects", projectID, "codes", codeID), func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "DELETE", r.Method) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"msg":"Deleted"}`)
	})
	ok, resp, err := client.Codes.DeleteCode(codeID)
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

func TestCodesServices_DockerLogin(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	muxIRON.HandleFunc(client.Path("projects", projectID, "credentials"), func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "POST", r.Method) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"msg":"Credentials added."}`)
	})
	ok, resp, err := client.Codes.DockerLogin(iron.DockerCredentials{
		Username:      "ron",
		Password:      "swanson",
		Email:         "ron.swanson@pawnee.gov",
		ServerAddress: "docker.io",
	})
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
