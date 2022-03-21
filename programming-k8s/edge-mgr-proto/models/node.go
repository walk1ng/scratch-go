package models

import (
	"edge-mgr-proto/types"

	corev1 "k8s.io/api/core/v1"
)

type NodeDetails struct {
	Node     *corev1.Node          `json:"node"`
	NodeInfo corev1.NodeSystemInfo `json:"nodeInfo"`
}

type ClusterNodesResponse struct {
	Count uint                   `json:"count"`
	Nodes []*ClusterNodeOverview `json:"nodes"`
}

type ClusterNodeOverview struct {
	Name                  string                    `json:"nodeName"`
	Hostname              string                    `json:"hostname"`
	InternalIP            string                    `json:"internalIP"`
	Status                types.NodeConditionStatus `json:"status"`
	PodCount              uint                      `json:"podCount"`
	CpuUsage              float64                   `json:"cpuUsage"`
	MemoryUsage           float64                   `json:"memoryUsage"`
	DiskUsage             float64                   `json:"diskUsage"`
	DiskIOUsage           float64                   `json:"diskIOUsage"`
	corev1.NodeSystemInfo `json:",inline"`
}
