package mdm_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/connect/mdm"
	"github.com/stretchr/testify/assert"
)

// TODO replace with BlobDataContract capture
func TestBlobDataContractCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	id := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	name := "TestContract"
	createdResource := `{
  "meta": {
    "lastUpdated": "2021-11-10T18:47:24.059503+00:00",
    "versionId": "ae6fabc0-3a65-4bb2-a08f-64d4e6453f0d"
  },
  "id": "` + id + `",
  "resourceType": "BlobDataContract",
  "name": "` + name + `",
  "dataTypeId": {
    "reference": "DataType/7b26ddb7-910b-4faf-b122-e1fd27356b14"
  },
  "rootPathInBucket": "foo",
  "bucketId": {
    "reference": "Bucket/8b26ddb7-910b-4faf-b122-e1fd27356b14"
  },
  "storageClassId": {
    "reference": "StorageClass/2b26ddb7-910b-4faf-b122-e1fd27356b14"
  }
}`
	muxMDM.HandleFunc("/connect/mdm/BlobDataContract", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, createdResource)
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "resourceType": "Bundle",
  "type": "searchset",
  "pageTotal": 0,
  "link": [
    {
      "relation": "string",
      "url": "string"
    }
  ],
  "entry": [
    {
      "fullUrl": "string",
      "resource": `+createdResource+`
    }
  ]
}`)
		}
	})
	muxMDM.HandleFunc("/connect/mdm/BlobDataContract/"+id, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, createdResource)
		case http.MethodPut:
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, createdResource)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})

	var c mdm.BlobDataContract
	c.Name = name
	c.RootPathInBucket = "foo"

	created, resp, err := mdmClient.BlobDataContracts.Create(c)
	if !assert.Nilf(t, err, "unexpected error: %v", err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, created) {
		return
	}
	created, resp, err = mdmClient.BlobDataContracts.GetByID(created.ID)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	if !assert.NotNil(t, createdResource) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, id, created.ID)

	ok, resp, err := mdmClient.BlobDataContracts.Delete(*created)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, createdResource)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode())
}
