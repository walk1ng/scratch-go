package types

type ClusterStatus string

const (
	ClusterStatusReady    ClusterStatus = "ready"
	ClusterStatusNotReady ClusterStatus = "not_ready"
	ClusterStatusUnknown  ClusterStatus = "unknown"
)

type NodeConditionStatus string

const (
	NodeStatusReady    NodeConditionStatus = "ready"
	NodeStatusNotReady NodeConditionStatus = "not_ready"
	NodeStatusUnknown  NodeConditionStatus = "unknown"
	NodeStatusNotExist NodeConditionStatus = "not_exist"
)
