package client

import (
	"fmt"
	"time"
)

type CNPGClusterPlan string

const (
	CNPGClusterPlanStarter    CNPGClusterPlan = "Starter"
	CNPGClusterPlanEnterprise CNPGClusterPlan = "Enterprise"
)

type ClusterStatus string

const (
	// CNPGClusterStatusCreating is the status of the cluster when it is creating.
	CNPGClusterStatusInitializing ClusterStatus = "Initializing"
	CNPGClusterStatusReady        ClusterStatus = "Ready"
	CNPGClusterStatusNotReady     ClusterStatus = "NotReady"
	CNPGClusterStatusDeleted      ClusterStatus = "Deleted"
	CNPGClusterStatusUpgrading    ClusterStatus = "Upgrading"
	CNPGClusterStatusSuspended    ClusterStatus = "Suspended"
	CNPGClusterStatusResuming     ClusterStatus = "Resuming"
)

type Region string

var (
	USEast1Region Region = "us-east-1"
	EUWest1Region Region = "eu-west-1"
)

var (
	ServerResourceAWST3XLarge  ServerResource = "aws-t3-xlarge-4c-16g"
	ServerResourceAWSM7ILarge  ServerResource = "aws-m7i-large-2c-8g"
	ServerResourceAWSR7ILarge  ServerResource = "aws-r7i-large-2c-16g"
	ServerResourceAWSR7IXLarge ServerResource = "aws-r7i-xlarge-4c-32g"
)

type CNPGClusterType string

const (
	// CNPGClusterTypeShared is the type of the cluster when it is shared.
	CNPGClusterTypeShared CNPGClusterType = "Shared"
	// CNPGClusterTypeDedicated is the type of the cluster when it is dedicated.
	CNPGClusterTypeDedicated CNPGClusterType = "Dedicated"
)

type ServerResource string

type ClusterProviderType string

const (
	AWSCloudProvider ClusterProviderType = "aws"
)

type ClusterProvider struct {
	Type   ClusterProviderType `json:"type,omitempty"`
	Region string              `json:"region,omitempty"`
}
type CNPGClusterPlanInfo struct {
	// Plan is the plan of the cluster.
	Plan        CNPGClusterPlan `json:"plan"`
	Description string          `json:"description"`
	Type        CNPGClusterType `json:"type"`
}

type CNPGCluster struct {
	Spec   CNPGClusterSpec   `json:"spec"`
	Status CNPGClusterStatus `json:"status"`
}

type CNPGClusterStatus struct {
	// Status is the status of the cluster.
	Status    ClusterStatus `json:"status,omitempty"`
	Endpoint  Endpoint      `json:"endpoint,omitempty"`
	ClusterID string        `json:"cluster_id,omitempty"`
	ProjectID string        `json:"project_id,omitempty"`
	UpdatedAt time.Time     `json:"updated_at,omitempty"`
}

type Endpoint struct {
	SuperUserEndpoint       string `json:"superuser_endpoint,omitempty"`
	VectorUserEndpoint      string `json:"vector_endpoint,omitempty"`
	PoolerSuperUserEndpoint string `json:"pooler_super_user_endpoint,omitempty"`
	PoolerUserEndpoint      string `json:"pooler_user_endpoint,omitempty"`
}

type CNPGClusterSpec struct {
	// ID is the id of the cluster.
	ID string `json:"id"`
	// Name is the name of the cluster.
	Name string `json:"name"`
	// ClusterProvider is the cluster provider of the cluster.
	ClusterProvider ClusterProvider `json:"cluster_provider"`
	// ServerResource is the server resource of the cluster instance.
	ServerResource   ServerResource   `json:"server_resource,omitempty"`
	PostgreSQLConfig PostgreSQLConfig `json:"postgresql_config,omitempty"`
	Plan             CNPGClusterPlan  `json:"plan"`
}

type CNPGClusterList struct {
	// Items is the list of the clusters.
	Items []CNPGCluster `json:"items"`
}

type PostgreSQLConfig struct {
	Instances      int          `json:"instances"`
	Image          string       `json:"image"`
	PGDataDiskSize string       `json:"pg_data_disk_size"`
	VectorConfig   VectorConfig `json:"vector_config,omitempty"`
	EnablePooler   bool         `json:"enable_pooler"`
}

type VectorConfig struct {
	DatabaseName string `json:"database_name"`
}

type CNPGClusterUpgradeRequest struct {
	// Plan is the plan of the cluster.
	Plan CNPGClusterPlan `json:"plan"`
	// ServerResource is the server resource of the cluster instance.
	ServerResource ServerResource `json:"server_resource,omitempty"`
	// PGDataDiskSize is the disk size of the Postgres PGData.
	PGDataDiskSize string `json:"pg_data_disk_size"`
}

func (c *Client) CreateCluster(params CNPGClusterSpec, userID string) (*CNPGCluster, error) {
	if params.PostgreSQLConfig.Instances == 0 {
		params.PostgreSQLConfig.Instances = 1
	}

	var clusterResponse CNPGCluster
	err := c.do("POST", fmt.Sprintf("users/%s/cnpgs", userID), params, &clusterResponse)
	return &clusterResponse, err
}

func (c *Client) GetCluster(userID string, clusterID string) (*CNPGCluster, error) {
	var clusterResponse CNPGCluster
	err := c.do("GET", fmt.Sprintf("users/%s/cnpgs/%s", userID, clusterID), nil, &clusterResponse)
	return &clusterResponse, err
}

func (c *Client) UpgradeCluster(userID string, clusterID string, params CNPGClusterUpgradeRequest) (*CNPGCluster, error) {
	var clusterResponse CNPGCluster
	err := c.do("PUT", fmt.Sprintf("users/%s/cnpgs/%s/upgrade", userID, clusterID), params, &clusterResponse)
	return &clusterResponse, err
}

func (c *Client) DeleteCluster(userID, clusterID string) error {
	return c.do("DELETE", fmt.Sprintf("users/%s/cnpgs/%s", userID, clusterID), nil, nil)
}
