package blr_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/blr"
	"github.com/stretchr/testify/assert"
)

func blobBody(id, dataType, state string) string {
	return fmt.Sprintf(`{
  "resourceType": "Blob",
  "meta": {
    "lastUpdated": "2022-05-25T19:36:10Z",
    "versionId": "1"
  },
  "id": "%s",
  "dataType": "%s",
  "autoGenerateBlobPathName": true,
  "bucket": "bucket-exact-moose",
  "blobPath": "dae89cf0-888d-4a26-8c1d-578e97365efc/64e403e6-d215-457a-bf12-2a4f49038208/tf-exact-moose/9b6f1d8a-0967-42d8-9622-d30c877c9da4/2022/05/25",
  "blobName": "74f78808-10d2-4199-92cf-e112bd9ff4de-19_36_10.bin",
  "createdBy": "9b6f1d8a-0967-42d8-9622-d30c877c9da4",
  "dataAccessUrl": "https://pre-signed.upload.url.com/something?x-amz-server-side-encryption=AES256",
  "dataAccessUrlExpiry": "2022-05-25T19:41:10.677Z",
  "creation": "2022-05-25T19:36:10Z",
  "multipartEnabled": false,
  "uploadOnBehalf": false,
  "managingOrganization": "dae89cf0-888d-4a26-8c1d-578e97365efc",
  "propositionGuid": "64e403e6-d215-457a-bf12-2a4f49038208",
  "state": "%s"
}`, id, dataType, state)
}

func TestBlobCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	blobID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc1"
	dataType := "tf-exact-moose"
	muxBLR.HandleFunc("/connect/blobrepository/Blob", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.Header().Set("Etag", "1")
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, blobBody(blobID, dataType, "uploading"))
		}
	})
	muxBLR.HandleFunc("/connect/blobrepository/Blob/"+blobID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, blobBody(blobID, dataType, "uploading"))
		case "PUT":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, blobBody(blobID, dataType, "uploading"))
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		}
	})

	created, resp, err := blrClient.Blobs.Create(blr.Blob{
		DataType: dataType,
	})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, created) {
		return
	}
	assert.Equal(t, dataType, created.DataType)
	assert.Equal(t, blobID, created.ID)
	assert.NotNil(t, created.State)

	res, resp, err := blrClient.Blobs.Delete(*created)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.True(t, res)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
