package tdr

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDataItem(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	muxTDR.HandleFunc("/store/tdr/DataItem", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("organization") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
			"type": "searchset",
			"total": 1,
			"entry": [
			  {
				"fullUrl": "https://foo-bar.com/store/tdr/DataItem?organization=TDROrg&device=deviceSystem%7C10001AD+UC-351+PBT+Ci&_id=73f31d0a73020d151a98c4b856db90e6b7558bcddef299ca112e820eb64a3426",
				"resource": {
				  "id": "73f31d0a73020d511a98c4b856db90e6b7558bcddef992ca112e820eb64a3426",
				  "meta": {
					"versionId": "1",
					"lastUpdated": "2018-07-13T10:08:00.488Z"
				  },
				  "organization": "TDROrg",
				  "dataType": {
					"system": "systemString",
					"code": "codeString"
				  },
				  "timestamp": "1970-01-18T17:14:16.000Z",
				  "device": {
					"system": "deviceSystem",
					"value": "SOMEDEVICE"
				  },
				  "data": {
					"string": "Foo",
					"integer": 28,
					"boolean": true,
					"dateTime": "2018-07-13T10:08:00.488Z"
				  },
				  "creationTimestamp": "2018-07-13T10:08:00.488Z",
				  "resourceType": "DataItem"
				}
			  }
			],
			"_startAt": 0,
			"link": [
			  {
				"relation": "next",
				"url": "https://foo-bar.com/store/tdr/DataItem?dataType=systemString%7CcodeString&organization=TDROrg&_startAt=1"
			  }
			],
			"resourceType": "Bundle"
		  }`)
	})
	dataItems, resp, err := tdrClient.DataItems.GetDataItem(&GetDataItemOptions{
		Organization: String("TDROrg"),
	}, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, 1, len(dataItems))

	dataItems, resp, err = tdrClient.DataItems.GetDataItem(&GetDataItemOptions{
		Organization: String("TDROrg"),
	}, DataSearch(KeyValue{"data.foo": "bar"}))
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(dataItems))
}
