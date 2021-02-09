package stl_test

import (
	"context"
	"github.com/philips-software/go-hsdp-api/stl"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestGetAppResourceByID(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	muxSTL.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "data": {
    "applicationResource": {
      "id": 1,
      "deviceId": 53615,
      "name": "terraform.yml",
      "content": "YXBpVmVyc2lvbjogdjEKa2luZDogU2VjcmV0Cm1ldGFkYXRhOgogIG5hbWU6IHNlY3JldC1zYS1zYW1wbGUKICBhbm5vdGF0aW9uczoKICAgIGt1YmVybmV0ZXMuaW8vc2VydmljZS1hY2NvdW50Lm5hbWU6ICJzYS1uYW1lIgp0eXBlOiBrdWJlcm5ldGVzLmlvL3NlcnZpY2UtYWNjb3VudC10b2tlbgpkYXRhOgogICMgWW91IGNhbiBpbmNsdWRlIGFkZGl0aW9uYWwga2V5IHZhbHVlIHBhaXJzIGFzIHlvdSBkbyB3aXRoIE9wYXF1ZSBTZWNyZXRzCiAgZXh0cmE6IFltRnlDZz09Cg=="
    }
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	ctx := context.Background()
	app, err := client.Apps.GetAppResourceByID(ctx, 1)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, app) {
		return
	}
	assert.Equal(t, int64(1), app.ID)
	assert.Equal(t, "terraform.yml", app.Name)
}

func TestGetAppResourceByDeviceIDAndName(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	muxSTL.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "data": {
    "applicationResource": {
      "id": 1,
      "deviceId": 53615,
      "name": "terraform.yml",
      "content": "YXBpVmVyc2lvbjogdjEKa2luZDogU2VjcmV0Cm1ldGFkYXRhOgogIG5hbWU6IHNlY3JldC1zYS1zYW1wbGUKICBhbm5vdGF0aW9uczoKICAgIGt1YmVybmV0ZXMuaW8vc2VydmljZS1hY2NvdW50Lm5hbWU6ICJzYS1uYW1lIgp0eXBlOiBrdWJlcm5ldGVzLmlvL3NlcnZpY2UtYWNjb3VudC10b2tlbgpkYXRhOgogICMgWW91IGNhbiBpbmNsdWRlIGFkZGl0aW9uYWwga2V5IHZhbHVlIHBhaXJzIGFzIHlvdSBkbyB3aXRoIE9wYXF1ZSBTZWNyZXRzCiAgZXh0cmE6IFltRnlDZz09Cg=="
    }
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	ctx := context.Background()
	app, err := client.Apps.GetAppResourceByDeviceIDAndName(ctx, 1, "terraform.yml")
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, app) {
		return
	}
	assert.Equal(t, int64(1), app.ID)
	assert.Equal(t, "terraform.yml", app.Name)
}

func TestGetAppResourcesBySerial(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	muxSTL.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "data": {
    "applicationResources": {
      "edges": [
        {
          "node": {
            "id": 874,
            "deviceId": 53615,
            "name": "ingress.yml",
            "content": "LS0tCmFwaVZlcnNpb246IGV4dGVuc2lvbnMvdjFiZXRhMQpraW5kOiBJbmdyZXNzCm1ldGFkYXRhOgogIG5hbWU6IGdvLWhlbGxvLXdvcmxkCiAgYW5ub3RhdGlvbnM6CiAgICBpbmdyZXNzLmt1YmVybmV0ZXMuaW8vc3NsLXJlZGlyZWN0OiAidHJ1ZSIKICAgIGt1YmVybmV0ZXMuaW8vaW5ncmVzcy5jbGFzczogdHJhZWZpawpzcGVjOgogIHJ1bGVzOgogIC0gaG9zdDogZGVtby5lZGdlLmhzZHAuaW8KICAgIGh0dHA6CiAgICAgIHBhdGhzOgogICAgICAtIHBhdGg6IC8KICAgICAgICBiYWNrZW5kOgogICAgICAgICAgc2VydmljZU5hbWU6IGdvLWhlbGxvLXdvcmxkLXNlcnZpY2UKICAgICAgICAgIHNlcnZpY2VQb3J0OiA4MA=="
          }
        },
        {
          "node": {
            "id": 876,
            "deviceId": 53615,
            "name": "service.yml",
            "content": "YXBpVmVyc2lvbjogdjEKa2luZDogU2VydmljZQptZXRhZGF0YToKICBuYW1lOiBnby1oZWxsby13b3JsZC1zZXJ2aWNlCiAgbGFiZWxzOgogICAgYXBwLmt1YmVybmV0ZXMuaW8vbmFtZTogZ28taGVsbG8td29ybGQtc2VydmljZQogICAgYXBwLmt1YmVybmV0ZXMuaW8vaW5zdGFuY2U6IGluaXRpYWwKICAgIGFwcC5rdWJlcm5ldGVzLmlvL3ZlcnNpb246ICIwLjAuMSIKc3BlYzoKICBwb3J0czoKICAtIHBvcnQ6IDgwCiAgICB0YXJnZXRQb3J0OiA4MDgwCiAgICBwcm90b2NvbDogVENQCiAgc2VsZWN0b3I6CiAgICBhcHAua3ViZXJuZXRlcy5pby9uYW1lOiBnby1oZWxsby13b3JsZA=="
          }
        },
        {
          "node": {
            "id": 877,
            "deviceId": 53615,
            "name": "deployment.yml",
            "content": "YXBpVmVyc2lvbjogYXBwcy92MQpraW5kOiBEZXBsb3ltZW50Cm1ldGFkYXRhOgogIG5hbWU6IGdvLWhlbGxvLXdvcmxkCiAgbGFiZWxzOgogICAgYXBwLmt1YmVybmV0ZXMuaW8vbmFtZTogZ28taGVsbG8td29ybGQKICAgIGFwcC5rdWJlcm5ldGVzLmlvL3ZlcnNpb246ICIwLjAuMSIKICAgIGFwcC5rdWJlcm5ldGVzLmlvL21hbmFnZWQtYnk6IENvbnNvbGUKc3BlYzoKICByZXBsaWNhczogMQogIHNlbGVjdG9yOgogICAgbWF0Y2hMYWJlbHM6CiAgICAgIGFwcC5rdWJlcm5ldGVzLmlvL25hbWU6IGdvLWhlbGxvLXdvcmxkCiAgdGVtcGxhdGU6CiAgICBtZXRhZGF0YToKICAgICAgbGFiZWxzOgogICAgICAgIGFwcC5rdWJlcm5ldGVzLmlvL25hbWU6IGdvLWhlbGxvLXdvcmxkCiAgICBzcGVjOgogICAgICBjb250YWluZXJzOgogICAgICAgIC0gbmFtZTogZ28taGVsbG8td29ybGQKICAgICAgICAgIHNlY3VyaXR5Q29udGV4dDoKICAgICAgICAgICAge30KICAgICAgICAgIGltYWdlOiAibG9hZm9lL2dvLWhlbGxvLXdvcmxkOmxhdGVzdCIKICAgICAgICAgIGltYWdlUHVsbFBvbGljeTogSWZOb3RQcmVzZW50CiAgICAgICAgICBlbnY6CiAgICAgICAgICBwb3J0czoKICAgICAgICAgICAgLSBuYW1lOiBodHRwCiAgICAgICAgICAgICAgY29udGFpbmVyUG9ydDogODA4MAogICAgICAgICAgICAgIHByb3RvY29sOiBUQ1AKICAgICAgICAgIGxpdmVuZXNzUHJvYmU6CiAgICAgICAgICAgIGh0dHBHZXQ6CiAgICAgICAgICAgICAgcGF0aDogLwogICAgICAgICAgICAgIHBvcnQ6IDgwODAKICAgICAgICAgIHJlYWRpbmVzc1Byb2JlOgogICAgICAgICAgICBodHRwR2V0OgogICAgICAgICAgICAgIHBhdGg6IC8KICAgICAgICAgICAgICBwb3J0OiA4MDgwCiAgICAgICAgICByZXNvdXJjZXM6CiAgICAgICAgICAgIHt9"
          }
        }
      ]
    }
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	ctx := context.Background()
	resources, err := client.Apps.GetAppResourcesBySerial(ctx, "serial")
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resources) {
		return
	}
	if !assert.Equal(t, 3, len(*resources)) {
		return
	}
	assert.Equal(t, "ingress.yml", (*resources)[0].Name)
	assert.Equal(t, "service.yml", (*resources)[1].Name)
	assert.Equal(t, "deployment.yml", (*resources)[2].Name)
}

func TestCreateAppResource(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	muxSTL.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "data": {
    "createApplicationResource": {
      "success": true,
      "message": "Successfully accepted create",
      "statusCode": 202,
      "requestId": "k3s-f4f57692-1674-417c-b1f7-01437091523f",
      "applicationResource": {
        "id": 1910,
        "deviceId": 53615,
        "name": "terraform.yml",
        "content": "YXBpVmVyc2lvbjogdjEKa2luZDogU2VjcmV0Cm1ldGFkYXRhOgogIG5hbWU6IHNlY3JldC1zYS1zYW1wbGUKICBhbm5vdGF0aW9uczoKICAgIGt1YmVybmV0ZXMuaW8vc2VydmljZS1hY2NvdW50Lm5hbWU6ICJzYS1uYW1lIgp0eXBlOiBrdWJlcm5ldGVzLmlvL3NlcnZpY2UtYWNjb3VudC10b2tlbgpkYXRhOgogICMgWW91IGNhbiBpbmNsdWRlIGFkZGl0aW9uYWwga2V5IHZhbHVlIHBhaXJzIGFzIHlvdSBkbyB3aXRoIE9wYXF1ZSBTZWNyZXRzCiAgZXh0cmE6IFltRnlDZz09Cg=="
      }
    }
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	serial := "foo"
	ctx := context.Background()
	app, err := client.Apps.CreateAppResource(ctx, stl.CreateApplicationResourceInput{
		SerialNumber: serial,
		Name:         "terraform.yml",
		Content:      "BLA",
	})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, app) {
		return
	}
	assert.Equal(t, "terraform.yml", app.Name)
}
