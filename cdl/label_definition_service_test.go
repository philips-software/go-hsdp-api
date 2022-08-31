package cdl_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/cdl"
	"github.com/stretchr/testify/assert"
)

func TestLabelDefinitionCRD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	studyID := "456aa456-ba82-485e-ac5d-014b4df33c4e"
	labelDefID := "f23843d0-e654-4033-9a35-8b9c497eaa66"

	labelDef := cdl.LabelDefinition{
		LabelDefName: "labeldefname",
		Description:  "desc",
		LabelScope: cdl.LabelScope{
			Type: "DataObject.DICOM",
		},
		Label: "label4",
		Type:  "cdl/video-classification",
		Labels: []cdl.LabelsArrayElem{
			{
				"labelarrayelem1",
			},
			{
				"labelarrayelem2",
			},
		},
	}

	muxCDL.HandleFunc("/store/cdl/"+cdlTenantID+"/Study/"+studyID+"/LabelDef", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{
    "id": "f23843d0-e654-4033-9a35-8b9c497eaa66",
    "labelDefName": "TF2",
    "description": "TF TEST",
    "labelScope": {
        "type": "DataObject.DICOM"
    },
    "researchStudyId": "a1467792-ef81-11eb-8ac2-477a9e3b09aa",
    "label": "videoQualityTF4",
    "type": "cdl/video-classification",
    "labels": [
        {
            "label": "good"
        },
        {
            "label": "bad"
        },
        {
            "label": "acceptable"
        },
        {
            "label": "something"
        },
        {
            "label": "something1"
        }
    ],
    "createdBy": "user@philips.com",
    "createdOn": "2021-07-28T17:43:33.488+00:00"
}`)
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
    "resourceType": "Bundle",
    "id": "b6d0c8e6-6df4-42ae-a239-898e92c5814d",
    "type": "searchset",
    "link": [
        {
            "relation": "self",
            "url": "https://cicd-datalake.cloud.pcftest.com/store/cdl/1f5be763-f896-4883-80fa-5593cd69556d/Study/a1467792-ef81-11eb-8ac2-477a9e3b09aa/LabelDef"
        }
    ],
    "entry": [
  {
    "fullUrl": "https://cicd-datalake.cloud.pcftest.com/store/cdl/1f5be763-f896-4883-80fa-5593cd69556d/Study/a1467792-ef81-11eb-8ac2-477a9e3b09aa/LabelDef/6429102c-a133-4326-85ba-12e5ff1ba7d3",
    "resource": {
      "id": "f23843d0-e654-4033-9a35-8b9c497eaa66",
      "resourceType": "LabelDef",
      "labelDefName": "TF1",
      "description": "TF TEST",
      "labelScope": {
        "type": "DataObject.DICOM"
      },
      "researchStudyId": "a1467792-ef81-11eb-8ac2-477a9e3b09aa",
      "label": "videoQualityTF",
      "type": "cdl/video-classification",
      "labels": [
        {
          "label": "good"
        },
        {
          "label": "bad"
        },
        {
          "label": "acceptable"
        },
        {
          "label": "something"
        }
      ],
      "createdBy": "user@philips.com",
      "createdOn": "2021-07-28T08:59:33.960+00:00"
    }
  }
]}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	muxCDL.HandleFunc("/store/cdl/"+cdlTenantID+"/Study/"+studyID+"/LabelDef/"+labelDefID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
      "id": "f23843d0-e654-4033-9a35-8b9c497eaa66",
      "resourceType": "LabelDef",
      "labelDefName": "TF1",
      "description": "TF TEST",
      "labelScope": {
        "type": "DataObject.DICOM"
      },
      "researchStudyId": "a1467792-ef81-11eb-8ac2-477a9e3b09aa",
      "label": "videoQualityTF",
      "type": "cdl/video-classification",
      "labels": [
        {
          "label": "good"
        },
        {
          "label": "bad"
        },
        {
          "label": "acceptable"
        },
        {
          "label": "something"
        }
      ],
      "createdBy": "user@philips.com",
      "createdOn": "2021-07-28T08:59:33.960+00:00"
    }
  }`)
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	createdLabelDef, resp, err := cdlClient.LabelDefinition.CreateLabelDefinition(studyID, labelDef)

	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, createdLabelDef) {
		return
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode())
	assert.Equal(t, createdLabelDef.ID, labelDefID)

	item, resp, err := cdlClient.LabelDefinition.GetLabelDefinitionByID(studyID, labelDefID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, item) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, labelDefID, item.ID)

	items, resp, err := cdlClient.LabelDefinition.GetLabelDefinitions(studyID, &cdl.GetOptions{})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, item) {
		return
	}
	assert.Equal(t, labelDefID, items[0].ID)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	resp, err = cdlClient.LabelDefinition.DeleteLabelDefinitionById(studyID, labelDefID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusNoContent, resp.StatusCode())

}
