package blr_test

import (
	"fmt"
	"github.com/philips-software/go-hsdp-api/blr"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func bucketBody(id, name string) string {
	return fmt.Sprintf(`{
  "resourceType": "Bucket",
  "id": "%s",
  "name": "%s",
  "enableHSDPDomain": false,
  "enableCDN": false,
  "priceClass": "ALL",
  "cacheControlAge": 0,
  "propositionId": {
    "reference": "string",
    "display": "string"
  },
  "corsConfiguration": {
    "allowedOrigins": [
      "string"
    ],
    "allowedMethods": [
      "PUT"
    ],
    "allowedHeaders": [
      "string"
    ],
    "exposeHeaders": [
      "string"
    ],
    "maxAgeSeconds": 1
  },
  "enableCreateOrDeleteBlobMeta": true
}`, id, name)
}

func bundleResponseBody(id, effect, action, principal, resource string) string {
	return fmt.Sprintf(`{
  "resourceType": "Bundle",
  "type": "searchset",
  "link": [
    {
      "relation": "self",
      "url": "BlobStorePolicy?_id=%s&_count=100"
    }
  ],
  "entry": [
    {
      "resource": {
        "resourceType": "BlobStorePolicy",
        "id": "%s",
        "statement": [
          {
            "principal": [
              "%s"
            ],
            "action": [
              "%s"
            ],
            "resource": [
              "%s"
            ],
            "effect": "%s"
          }
        ]
      },
      "fullUrl": "https://foo.bar.com/connect/blobrepository/configuration/BlobStorePolicy/%s"
    }
  ]
}`, id, id, principal, action, effect, resource, id)
}

func blobStorePolicyBody(id, effect, action, principal, resource string) string {
	return fmt.Sprintf(`{
  "resourceType": "BlobStorePolicy",
  "meta": {
    "lastUpdated": "2022-05-25T19:36:10Z",
    "versionId": "1"
  },
  "id": "%s",
  "statement": [
    {
      "effect": "%s",
      "action": [
        "%s"
      ],
      "principal": [
        "%s"
      ],
      "resource": [
        "%s"
      ]
    }
  ]
}`, id, effect, action, principal, resource)
}

func TestBlobStorePolicyCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	blobStorePolicyID := "dbf1d779-ab9f-4c27-b4aa-ea75f9efbbc1"
	muxBLR.HandleFunc("/connect/blobrepository/configuration/BlobStorePolicy", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.Header().Set("Etag", "1")
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, blobStorePolicyBody(blobStorePolicyID, "effect", "action", "principal", "resource"))
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, bundleResponseBody(blobStorePolicyID, "effect", "action", "principal", "resource"))
		case "DELETE", "PUT":
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	muxBLR.HandleFunc("/connect/blobrepository/configuration/BlobStorePolicy/"+blobStorePolicyID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusMethodNotAllowed)
		case "PUT":
			w.WriteHeader(http.StatusMethodNotAllowed)
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		}
	})

	created, resp, err := blrClient.Configurations.CreateBlobStorePolicy(blr.BlobStorePolicy{
		Statement: []blr.BlobStorePolicyStatement{
			{
				Effect:    "effect",
				Action:    []string{"action"},
				Principal: []string{"principal"},
				Resource:  []string{"resource"},
			},
		},
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
	assert.Equal(t, []string{"action"}, created.Statement[0].Action)
	assert.Equal(t, blobStorePolicyID, created.ID)

	res, resp, err := blrClient.Configurations.DeleteBlobStorePolicy(*created)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.True(t, res)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode())

	found, resp, err := blrClient.Configurations.GetBlobStorePolicyByID(blobStorePolicyID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, found) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, blobStorePolicyID, found.ID)

}

func TestBucketCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	bucketID := "d3898a8e-8ef7-48db-9eb6-3f423e85a853"
	propositionID := "f6022cce-afd9-456f-bc02-1a423fe1e16a"
	name := "expanse"
	muxBLR.HandleFunc("/connect/blobrepository/configuration/Bucket", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.Header().Set("Etag", "1")
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, bucketBody(bucketID, name))
		}
	})
	muxBLR.HandleFunc("/connect/blobrepository/configuration/Bucket/"+bucketID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, bucketBody(bucketID, name))
		case "PUT":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, bucketBody(bucketID, name))
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		}
	})

	created, resp, err := blrClient.Configurations.CreateBucket(blr.Bucket{
		Name: name,
		PropositionID: blr.Reference{
			Reference: propositionID,
		},
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
	assert.Equal(t, bucketID, created.ID)

	res, resp, err := blrClient.Configurations.DeleteBucket(*created)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.True(t, res)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode())

	found, resp, err := blrClient.Configurations.GetBucketByID(bucketID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, found) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, bucketID, found.ID)

}
