package models

type Cluster struct {
	Name        string `json:"name,omitempty" binding:"required"`
	Description string `json:"description,omitempty"`
	KubeConfig  string `json:"kubeconfig,omitempty" binding:"required"`
	NodeCount   int    `json:"nodeCount,omitempty"`
	Version     string `json:"version,omitempty"`
}

type K8sClusterSummary struct{}

type Node struct {
}

type NodeSummary struct{}
