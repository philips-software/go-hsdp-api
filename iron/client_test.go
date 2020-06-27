package iron_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/philips-software/go-hsdp-api/iron"

	"github.com/stretchr/testify/assert"
)

var (
	muxIRON    *http.ServeMux
	serverIRON *httptest.Server
	client     *iron.Client
	projectID  = "48a0183d-a588-41c2-9979-737d15e9e860"
	token      = "YM7eZakYwqoui5znoH4g"
)

func setup(t *testing.T) func() {
	muxIRON = http.NewServeMux()
	serverIRON = httptest.NewServer(muxIRON)

	var err error

	client, err = iron.NewClient(&iron.Config{
		BaseURL:   serverIRON.URL,
		ProjectID: projectID,
		Token:     token,
		Debug:     true,
		DebugLog:  "/tmp/iron_test.log",
	})
	assert.Nil(t, err)
	assert.NotNil(t, client)

	return func() {
		serverIRON.Close()
	}
}

func TestClient_Error(t *testing.T) {
	muxIRON = http.NewServeMux()
	serverIRON = httptest.NewServer(muxIRON)

	var err error
	client, err = iron.NewClient(&iron.Config{
		BaseURL:   serverIRON.URL,
		ProjectID: projectID,
		Token:     token,
		Debug:     true,
	})
	assert.Nil(t, err)
	assert.NotNil(t, client)

	task, resp, err := client.Tasks.GetTasks()
	assert.NotNil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.NotNil(t, task)
	client.Close()
}

func TestClusterInfo_Encrypt(t *testing.T) {
	pubkey := []byte("-----BEGIN PUBLIC KEY----- MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCdS2oE9+dhexZc3/sEtI+a6ZKt 6FwBZaAgytdkQ7sX4FwbZAdJ7zFS1m0gDezyFTBJSPVjYOKYr0fu1ao/xkNkKnnz J2WkW6qsDNKwJgrHiCO1asnoW5XWtk8Yc4kKkg63REuV20x+QoD6onTCo3T2DfUI vZ8QOSJQ7NotGuO2wwIDAQAB -----END PUBLIC KEY-----")
	ci := &iron.ClusterInfo{
		UserID:      "foo",
		ClusterID:   "bar",
		ClusterName: "cluster",
	}
	encrypted, err := ci.Encrypt([]byte(`hello world`))
	assert.Equal(t, iron.ErrNoPublicKey, err)
	ci.Pubkey = string(pubkey)
	encrypted, err = ci.Encrypt([]byte(`hello world`))
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 224, len(encrypted))
}
