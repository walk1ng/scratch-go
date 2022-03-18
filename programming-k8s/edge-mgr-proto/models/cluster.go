package models

import (
	"edge-mgr-proto/types"
)

type Cluster struct {
	Name        string `json:"name,omitempty" binding:"required"`
	Description string `json:"description,omitempty"`
	KubeConfig  string `json:"kubeconfig,omitempty" binding:"required"`
	NodeCount   int    `json:"nodeCount,omitempty"`
	Version     string `json:"version,omitempty"`
}

type ClusterNodesRequest struct {
	Nodes []string `json:"nodes" binding:"required"`
}

type ClusterOverviewResponse struct {
	KubernetesVersion string              `json:"kubernetesVersion"`
	Status            types.ClusterStatus `json:"status"`
	PodCapacity       uint                `json:"podCapacity"`
	PodCount          uint                `json:"podCount"`
	NodeCount         uint                `json:"nodeCount"`
	Nodes             []string            `json:"nodes"`
	MasterCount       uint                `json:"masterCount"`
	Masters           []string            `json:"masters"`
	WorkerCount       uint                `json:"workerCount"`
	Workers           []string            `json:"workers"`
	CpuTotal          float64             `json:"cpuTotal"`
	CpuUsed           float64             `json:"cpuUsed"`
	MemoryBytesTotal  float64             `json:"memoryBytesTotal"`
	MemoryBytesUsed   float64             `json:"memoryBytesUsed"`
	DiskBytesTotal    float64             `json:"diskBytesTotal"`
	DiskBytesUsed     float64             `json:"diskBytesUsed"`
}

type ClusterNodesResponse struct {
	Count uint                   `json:"count"`
	Nodes []*ClusterNodeOverview `json:"nodes"`
}

type ClusterNodeOverview struct {
	Name        string                    `json:"hostnameOrIP"`
	Status      types.NodeConditionStatus `json:"status"`
	PodCount    uint                      `json:"podCount"`
	CpuUsage    float64                   `json:"cpuUsage"`
	MemoryUsage float64                   `json:"memoryUsage"`
	DiskUsage   float64                   `json:"diskUsage"`
	DiskIOUsage float64                   `json:"diskIOUsage"`
}
