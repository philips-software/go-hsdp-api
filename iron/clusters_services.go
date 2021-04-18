package iron

import (
	"time"
)

// ClustersServices implements API calls to get
// details on Iron clusters. As HSDP Iron clusters are not
// user serviceable only informational API calls are implemented
type ClustersServices struct {
	client    *Client
	projectID string
}

// Cluster describes an Iron cluster
type Cluster struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	UserID           string        `json:"user_id"`
	Memory           int64         `json:"memory"`
	DiskSpace        int64         `json:"disk_space"`
	CPUShare         int           `json:"cpu_share"`
	SharedTo         []ClusterUser `json:"shared_to"`
	RunnersTotal     int           `json:"runners_total"`
	RunnersAvailable int           `json:"runners_available"`
	Machines         []Machine     `json:"machines"`
}

// Machine is a node in an Iron cluster
type Machine struct {
	InstanceID       string `json:"instance_id"`
	Version          string `json:"version"`
	RunnersTotal     int    `json:"runners_total"`
	RunnersAvailable int    `json:"runners_available"`
}

// ClusterUser can share resources on a cluster
type ClusterUser struct {
	UserID    string    `json:"user_id"`
	ClusterID string    `json:"cluster_id"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
}

type ClusterStats struct {
	RunnersTotal     int       `json:"runners_total"`
	RunnersAvailable int       `json:"runners_available"`
	Queued           int       `json:"queued"`
	Ewma             int       `json:"ewma"`
	EwmaBeta         int       `json:"ewma_beta"`
	Instances        []Machine `json:"instances"`
}

// GetClusters gets the list of available clusters
// In some cases a token might not have the proper scope
// to retrieve a list of clusters in which case the list will be empty
func (c *ClustersServices) GetClusters() (*[]Cluster, *Response, error) {
	page := 0
	perPage := 100
	req, err := c.client.newRequest("GET", c.client.Path("clusters"), pageOptions{
		PerPage: &perPage,
		Page:    &page,
	}, nil)
	if err != nil {
		return nil, nil, err
	}
	var clusters struct {
		Clusters []Cluster `json:"clusters"`
	}
	resp, err := c.client.do(req, &clusters)
	return &clusters.Clusters, resp, err
}

// GetCluster gets cluster details
func (c *ClustersServices) GetCluster(clusterID string) (*Cluster, *Response, error) {
	req, err := c.client.newRequest("GET", c.client.Path("clusters", clusterID), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	var cluster struct {
		Cluster Cluster `json:"cluster"`
	}
	resp, err := c.client.do(req, &cluster)
	return &cluster.Cluster, resp, err
}

// GetClusterStats gets cluster statistics
func (c *ClustersServices) GetClusterStats(clusterID string) (*ClusterStats, *Response, error) {
	req, err := c.client.newRequest("GET", c.client.Path("clusters", clusterID, "stats"), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	var cluster struct {
		ClusterStats ClusterStats `json:"cluster"`
	}
	resp, err := c.client.do(req, &cluster)
	return &cluster.ClusterStats, resp, err
}
