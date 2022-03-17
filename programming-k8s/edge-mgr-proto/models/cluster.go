package models

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
	KubernetesVersion string   `json:"kubernetesVersion"`
	Status            string   `json:"status"`
	PodCapacity       uint     `json:"podCapacity"`
	Pods              uint     `json:"pods"`
	Nodes             uint     `json:"nodes"`
	NodeList          []string `json:"nodeList"`
	Masters           uint     `json:"masters"`
	MasterList        []string `json:"masterList"`
	Workers           uint     `json:"workers"`
	WorkerList        []string `json:"workerList"`
	CpuTotal          float64  `json:"cpuTotal"`
	CpuUsed           float64  `json:"cpuUsed"`
	MemoryBytesTotal  float64  `json:"memoryBytesTotal"`
	MemoryBytesUsed   float64  `json:"memoryBytesUsed"`
	DiskBytesTotal    float64  `json:"diskBytesTotal"`
	DiskBytesUsed     float64  `json:"diskBytesUsed"`
}

type ClusterNodesResponse struct {
	Count uint                   `json:"count"`
	Nodes []*ClusterNodeOverview `json:"nodes"`
}

type ClusterNodeOverview struct {
	Name             string  `json:"hostnameOrIP"`
	Status           string  `json:"status"`
	Pods             uint    `json:"pods"`
	CpuTotal         float64 `json:"cpuTotal"`
	CpuUsed          float64 `json:"cpuUsed"`
	MemoryBytesTotal float64 `json:"memoryBytesTotal"`
	MemoryBytesUsed  float64 `json:"memoryBytesUsed"`
	DiskBytesTotal   float64 `json:"diskBytesTotal"`
	DiskBytesUsed    float64 `json:"diskBytesUsed"`
}
