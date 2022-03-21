package types

// cluster status
type ClusterStatus string

const (
	ClusterStatusReady    ClusterStatus = "ready"
	ClusterStatusNotReady ClusterStatus = "not_ready"
	ClusterStatusUnknown  ClusterStatus = "unknown"
)

// node status
type NodeConditionStatus string

const (
	NodeStatusReady    NodeConditionStatus = "ready"
	NodeStatusNotReady NodeConditionStatus = "not_ready"
	NodeStatusUnknown  NodeConditionStatus = "unknown"
	NodeStatusNotExist NodeConditionStatus = "not_exist"
)

// node
type NodeMetricsTarget string

const (
	TargetNodeCPU    NodeMetricsTarget = "node_cpu"
	TargetNodeMemory NodeMetricsTarget = "node_memory"
	TargetNodeDisk   NodeMetricsTarget = "node_disk"
	TargetNodeDiskIO NodeMetricsTarget = "node_diskio"
)

type ClusterMetricsTarget string

const (
	TargetClusterUsedCPU     ClusterMetricsTarget = "cluster_cpu_used"
	TargetClusterTotalCPU    ClusterMetricsTarget = "cluster_cpu_total"
	TargetClusterUsedMemory  ClusterMetricsTarget = "cluster_memory_used"
	TargetClusterTotalMemory ClusterMetricsTarget = "cluster_memory_total"
	TargetClusterUsedDisk    ClusterMetricsTarget = "cluster_disk_used"
	TargetClusterTotalDisk   ClusterMetricsTarget = "cluster_disk_total"
)

type ClusterCpuUsage struct {
	Used  float64 `json:"used"`
	Total float64 `json:"total"`
}

type ClusterMemUsage struct {
	BytesUsed  float64 `json:"bytesUsed"`
	BytesTotal float64 `json:"bytesTotal"`
}

type ClusterDiskUsage struct {
	BytesUsed  float64 `json:"bytesUsed"`
	BytesTotal float64 `json:"bytesTotal"`
}
