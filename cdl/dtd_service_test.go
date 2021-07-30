package cdl_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/cdl"
	"github.com/stretchr/testify/assert"
)

func TestDataTypeDefinitionCRD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	dtdID := "123aa456-ba82-485e-ac5d-014b4df33c4e"
	updatedDtdDescription := "dtd for tf test updated"

	dataTypeDefToCreate := cdl.DataTypeDefinition{
		Name:        "dtdtestingonetwothree",
		Description: "dtd for tf test",
	}

	err := json.Unmarshal([]byte(`{"key": "value"}`), &dataTypeDefToCreate.JsonSchema)
	if err != nil {
		fmt.Println("err from json unmarshall of json schema", err)
	}

	muxCDL.HandleFunc("/store/cdl/"+cdlTenantID+"/DataTypeDefinition", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{
        "id": "123aa456-ba82-485e-ac5d-014b4df33c4e",
        "name": "dtdtestingonetwothree",
        "description": "dtd for tf test",
        "createdOn": "2021-07-25T08:18:17.441+00:00",
        "createdBy": "user@philips.com",
        "updatedOn": "2021-07-25T08:18:17.441+00:00",
        "updatedBy": "user@philips.com",
        "jsonSchema": {
            "key": "value"
        }
    }`)
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `[{
        "id": "123aa456-ba82-485e-ac5d-014b4df33c4e",
        "name": "dtdtestingonetwothree",
        "description": "dtd for tf test",
        "createdOn": "2021-07-25T08:18:17.441+00:00",
        "createdBy": "user@philips.com",
        "updatedOn": "2021-07-25T08:18:17.441+00:00",
        "updatedBy": "user@philips.com",
        "jsonSchema": {
            "key": "value"
        }
    }]`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	muxCDL.HandleFunc("/store/cdl/"+cdlTenantID+"/DataTypeDefinition/"+dtdID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
        "id": "123aa456-ba82-485e-ac5d-014b4df33c4e",
        "name": "dtdtestingonetwothree",
        "description": "dtd for tf test",
        "createdOn": "2021-07-25T08:18:17.441+00:00",
        "createdBy": "user@philips.com",
        "updatedOn": "2021-07-25T08:18:17.441+00:00",
        "updatedBy": "user@philips.com",
        "jsonSchema": {
            "key": "value"
        }
    }`)
		case "PUT":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
        "id": "123aa456-ba82-485e-ac5d-014b4df33c4e",
        "name": "dtdtestingonetwothree",
        "description": "dtd for tf test updated",
        "createdOn": "2021-07-25T08:18:17.441+00:00",
        "createdBy": "user@philips.com",
        "updatedOn": "2021-07-25T08:18:17.441+00:00",
        "updatedBy": "user@philips.com",
        "jsonSchema": {
            "key": "value",
			"updated_key": "value"
        }
    }`)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	created, resp, err := cdlClient.DataTypeDefinition.CreateDataTypeDefinition(dataTypeDefToCreate)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, created) {
		return
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, created.ID, dtdID)

	item, resp, err := cdlClient.DataTypeDefinition.GetDataTypeDefinitionByID(dtdID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, item) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, dtdID, item.ID)

	item, resp, err = cdlClient.DataTypeDefinition.UpdateDataTypeDefinition(*item)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, item) {
		return
	}
	assert.Equal(t, dtdID, item.ID)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, item.Description, updatedDtdDescription)

	listOfDtd, listOfDtdResp, err := cdlClient.DataTypeDefinition.GetDataTypeDefinitions(&cdl.GetOptions{})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, listOfDtdResp) {
		return
	}
	if !assert.NotNil(t, listOfDtd) {
		return
	}
	assert.Equal(t, http.StatusOK, listOfDtdResp.StatusCode)
	assert.Equal(t, len(listOfDtd), 1)
	assert.Equal(t, listOfDtd[0].ID, dtdID)
}
