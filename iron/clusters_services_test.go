package iron_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClustersServices_GetClusters(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	muxIRON.HandleFunc(client.Path("clusters"), func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "GET", r.Method) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"clusters": [] }`)
	})

	clusters, resp, err := client.Clusters.GetClusters()
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, clusters) {
		return
	}
	assert.Equal(t, 0, len(*clusters))
}

func TestClustersServices_GetCluster(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	clusterID := "Q3b9CZmGFEvTlr83RC4VUoxQ"
	muxIRON.HandleFunc(client.Path("clusters", clusterID), func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "GET", r.Method) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
  "cluster": {
    "id": "`+clusterID+`",
    "name": "dev_large_encrypted",
    "user_id": "FH959DA78Ygd6MejhhuQbNpq",
    "memory": 12884901888,
    "disk_space": 53687091200,
    "cpu_share": 2,
    "shared_to": [
      {
        "user_id": "iEhFUoe2O0UZvJP7Ezapb2Ja",
        "cluster_id": "`+clusterID+`",
        "updated_at": "2017-01-12T19:32:05.74Z",
        "name": "0fh8loao4z@pawnee.com",
        "email": "0fh8loao4z@pawnee.com"
      },
      {
        "user_id": "Prenp3fQGytw30UfJLuRfONo",
        "cluster_id": "`+clusterID+`",
        "updated_at": "2017-05-03T19:12:31.64Z",
        "name": "sb9cuyq74s@pawnee.com",
        "email": "sb9cuyq74s@pawnee.com"
      }
    ],
    "runners_total": 3,
    "runners_available": 2,
    "machines": [
      {
        "instance_id": "i-PaMqN6HVprovO0mcM",
        "version": "3.1.7-beta-2017-08-23.1",
        "runners_total": 3,
        "runners_available": 3
      },
      {
        "instance_id": "i-Xw2PJmNcRbhibkWNs",
        "version": "3.1.7-beta-2017-08-23.1",
        "runners_total": 3,
        "runners_available": 3
      },
      {
        "instance_id": "i-Gse1dURQbojPuDkon",
        "version": "3.1.7-beta-2017-08-23.1",
        "runners_total": 3,
        "runners_available": 2
      }
    ]
  }
}`)
	})

	cluster, resp, err := client.Clusters.GetCluster(clusterID)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, cluster) {
		return
	}
	assert.Equal(t, clusterID, cluster.ID)
}

func TestClustersServices_GetClusterStats(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	clusterID := "Q3b9CZmGFEvTlr83RC4VUoxQ"
	muxIRON.HandleFunc(client.Path("clusters", clusterID, "stats"), func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, "GET", r.Method) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
  "cluster": {
    "runners_total": 9,
    "runners_available": 6,
    "queued": 3,
    "ewma": 0,
    "ewma_beta": 0,
    "instances": [
      {
        "instance_id": "i-pW5ez5ZQkWkB4PblP",
        "version": "3.1.7-beta-2017-08-23.1",
        "runners_total": 3,
        "runners_available": 0
      },
      {
        "instance_id": "i-MqW4MjUe8rWg0cxpQ",
        "version": "3.1.7-beta-2017-08-23.1",
        "runners_total": 3,
        "runners_available": 3
      },
      {
        "instance_id": "i-2ndXr2ZRaBbL0cfxU",
        "version": "3.1.7-beta-2017-08-23.1",
        "runners_total": 3,
        "runners_available": 3
      }
    ]
  }
}`)
	})

	cluster, resp, err := client.Clusters.GetClusterStats(clusterID)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, cluster) {
		return
	}
	assert.Equal(t, 9, cluster.RunnersTotal)
	assert.Equal(t, 6, cluster.RunnersAvailable)

}
