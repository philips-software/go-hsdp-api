package tdr

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContract(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	if tdrClient == nil {
		t.Fatal("Expected tdrClient to be set")
	}

	muxTDR.HandleFunc("/store/tdr/Contract", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("dataType") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
			"type": "searchset",
			"total": 2,
			"entry": [
				{
					"fullUrl": "https://tdr-service-staging.us-east.philips-healthsuite.com/store/tdr/Contract?dataType=TestGo%7CTestGoContract1",
					"resource": {
						"deletePolicy": {
							"duration": 1,
							"unit": "MONTH"
						},
						"sendNotifications": false,
						"id": "TestGo|TestGoContract1",
						"meta": {
							"versionId": "1",
							"lastUpdated": "2018-11-05T12:55:23.700Z"
						},
						"organization": "TDROrg",
						"dataType": {
							"system": "TestGo",
							"code": "TestGoContract1"
						},
						"schema": {
							"$schema": "http://json-schema.org/draft-04/schema#",
							"type": "object",
							"properties": {
								"Temperature": {
									"type": "number"
								},
								"HeartRate": {
									"type": "integer"
								},
								"IsManualMeasurement": {
									"type": "boolean"
								},
								"DeviceStatus": {
									"type": "string"
								}
							},
							"required": [
								"Temperature",
								"HeartRate"
							]
						},
						"resourceType": "Contract"
					}
				},
				{
					"fullUrl": "https://tdr-service-staging.us-east.philips-healthsuite.com/store/tdr/Contract?dataType=TestGo%7CTestGoContract2",
					"resource": {
						"deletePolicy": {
							"duration": 1,
							"unit": "MONTH"
						},
						"sendNotifications": false,
						"id": "TestGo|TestGoContract2",
						"meta": {
							"versionId": "1",
							"lastUpdated": "2018-11-05T12:55:23.700Z"
						},
						"organization": "TDROrg",
						"dataType": {
							"system": "TestGo",
							"code": "TestGoContract2"
						},
						"schema": {
							"$schema": "http://json-schema.org/draft-04/schema#",
							"type": "object",
							"properties": {
								"Temperature": {
									"type": "number"
								},
								"HeartRate": {
									"type": "integer"
								},
								"IsManualMeasurement": {
									"type": "boolean"
								},
								"DeviceStatus": {
									"type": "string"
								}
							},
							"required": [
								"Temperature",
								"HeartRate"
							]
						},
						"resourceType": "Contract"
					}
				}
			],
			"resourceType": "Bundle"
		  }`)
	})
	contracts, resp, err := tdrClient.Contracts.GetContract(&GetContractOptions{
		DataType: String("TestGo|TestGoContract"),
	}, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, 2, len(contracts))
}

func TestCreateContract(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	if tdrClient == nil {
		t.Fatal("Expected tdrClient to be set")
	}

	muxTDR.HandleFunc("/store/tdr/Contract", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Unexpected EOF from reading request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var contract Contract
		err = json.Unmarshal(body, &contract)
		if err != nil {
			t.Errorf("Expected contract in body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", "https://golang-testurl.com/store/tdr/Contract?dataType=TestGo%7CTestGoContract")
		w.WriteHeader(http.StatusCreated)
	})

	var schemaContract = []byte(`{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"type": "object",
		"properties": {
		  "Temperature": {
			"type": "number"
		  },
		  "HeartRate": {
			"type": "integer"
		  },
		  "IsManualMeasurement": {
			"type": "boolean"
		  },
		  "DeviceStatus": {
			"type": "string"
		  }
		},
		"required": [
		  "Temperature",
		  "HeartRate"
		]
	  }`)

	var newContract = Contract{
		SendNotifications: false,
		Organization:      "DevOrg",
		DataType: DataType{
			System: "TestGo",
			Code:   "TestGoContract",
		},
		DeletePolicy: DeletePolicy{
			Duration: 1,
			Unit:     "MONTH",
		},
		Schema: json.RawMessage(schemaContract),
	}
	ok, resp, err := tdrClient.Contracts.CreateContract(newContract)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusCreated, resp.StatusCode())
	assert.Equal(t, true, ok, "expected contract creation to succeed")
}
