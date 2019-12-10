package credentials

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccess(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	productKey := "430deb9e-01c8-4a3b-81dd-e2e46569cd5e"

	muxCreds.HandleFunc("/core/credentials/Access", func(w http.ResponseWriter, r *http.Request) {
		if k := r.Header.Get("X-Product-Key"); k != productKey {
			t.Errorf(ErrMissingProductKey.Error())
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if k := r.Header.Get("Authorization"); k == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `[
				{
					"allowed": {
					  "resources": [
						""
					  ],
					  "actions": [
						"GET",
						"PUT",
						"DELETE"
					  ]
					},
					"credentials": {
					  "accessKey": "XAW9kNtN0klc5IxCdLLH",
					  "secretKey": "9zNj5Sju55wTXTaUla4tTXEZWbUKK7JXTeW8r8Tr",
					  "sessionToken": "r2Qe9JbZA3EpGqQcLGjYFhEUktKMYnmccWlb2wFRzCUsE393MNpQz6v5wsDSXc8VDZ78Apji0tJLNBMcnVHR6x93VDyra73ztfQUze9llsf7dm0dx6B398kdwEbrK62jd0whCWyk8tczkUei9dJpPediD4BbGRbDbP5w8W2GWoO3EHXS9p5qpPZu8seb1kRJ6o74fblxbc4UMnnq7csBwJxMuzc4iUU3T5SBMC2uRQmKRAYBfhVNoFpeiSxzqsNHBmbbBxQTNMy4vKCqHIhLXy9GRS2TE85TlVvDZF1eqjIVGeqtOHD18MUN5juLLatBT6uDPoJi10QgEUj7Zx2BbVIVud71jULt83EaJEkpU1fEghIgDgqgB3c5ZFNnxktWTR1fvfoALejkzqSaka7AKxnTZJaLsAeqBz5yTOl8N4pPNSBKpW8RkCIydHpq5NYj9qAH1QhEHJOSLAD1v4mzM6Rjofh2Q84PApA2oOXbJ31vyBLmXy8O7eAGQ98fzhnd3f3SCqTjqqksZiDZzfH1EQh83kYdsGDCTgTQJIRBRvLujvZU28hKyqLvN12fpPkZyraX5SXa4TN7m6vIfT8SjFEEZwTlafCj8dxa",
					  "expires": "2019-02-28T23:13",
					  "bucket": "de-ad-497838f9-9cb3-4b7b-8501-e5866b9a48e3"
					}
				  }
			]`)
		}
	})

	access, resp, err := credsClient.Access.GetAccess(&GetAccessOptions{ProductKey: &productKey})
	assert.Nil(t, err)
	assert.NotNil(t, access)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, len(access), "expected one access record")
	assert.Equal(t, "XAW9kNtN0klc5IxCdLLH", access[0].Credentials.AccessKey)
	assert.Equal(t, "9zNj5Sju55wTXTaUla4tTXEZWbUKK7JXTeW8r8Tr", access[0].Credentials.SecretKey)
}
