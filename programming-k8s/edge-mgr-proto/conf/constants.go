package conf

const (
	DISK_FSTYPE = "ext[234]|btrfs|xfs|zfs"

	/*
		Cluster metrics
	*/

	// cluster cpu
	QueryClusterCpuUsed  = `sum(irate(node_cpu_seconds_total{job="node-exporter", mode!="idle", instance=~"%s"}[5m]))`
	QueryClusterCpuCount = `sum(count without(cpu, mode) (node_cpu_seconds_total{job="node-exporter", mode="idle", instance=~"%s"}))`

	// cluster memory
	QueryClusterMemoryTotal = `sum(node_memory_MemTotal_bytes{job="node-exporter", instance=~"%s"})`
	QueryClusterMemoryUsed  = `sum(node_memory_MemTotal_bytes{job="node-exporter", instance=~"%s"}) - 
	sum(node_memory_MemFree_bytes{job="node-exporter", instance=~"%s"}) - 
	sum(node_memory_Buffers_bytes{job="node-exporter", instance=~"%s"}) - 
	sum(node_memory_Cached_bytes{job="node-exporter", instance=~"%s"}) + 
	sum(node_memory_Shmem_bytes{job="node-exporter", instance=~"%s"})`

	// cluster disk
	QueryClusterDiskTotal = `sum(node_filesystem_size_bytes{job="node-exporter", instance=~"%s", fstype!="", mountpoint="/"})`
	QueryClusterDiskUsed  = `sum(node_filesystem_size_bytes{job="node-exporter", instance=~"%s", fstype!="", mountpoint="/"}) -
	sum(node_filesystem_free_bytes{job="node-exporter", instance=~"%s", fstype!="", mountpoint="/"})`

	/*
		Node metrics
	*/
	// node cpu usage
	QueryNodeCpuUsage = `sum(irate(node_cpu_seconds_total{job="node-exporter", mode!="idle", instance="%s"}[3m])) / sum(count without(cpu, mode) (node_cpu_seconds_total{job="node-exporter", mode="idle", instance="%s"})) * 100`

	// node memory usage
	QueryNodeMemoryUsage = `(sum(node_memory_MemTotal_bytes{job="node-exporter", instance="%s"}) - 
	sum(node_memory_MemFree_bytes{job="node-exporter", instance="%s"}) - 
	sum(node_memory_Buffers_bytes{job="node-exporter", instance="%s"}) - 
	sum(node_memory_Cached_bytes{job="node-exporter", instance="%s"}) + 
	sum(node_memory_Shmem_bytes{job="node-exporter", instance="%s"})) / 
	sum(node_memory_MemTotal_bytes{job="node-exporter", instance="%s"}) * 100`

	// node disk usage
	QueryNodeDiskUsage = `(1 - 
		(max without (mountpoint, fstype) (node_filesystem_avail_bytes{job="node-exporter", fstype!="", instance="%s", mountpoint="/"})
		/ 
		max without (mountpoint, fstype) (node_filesystem_size_bytes{job="node-exporter", fstype!="", instance="%s", mountpoint="/"}))
		) * 100 `

	// node disk io usage
	QueryNodeDiskIOUsage = `max(rate(node_disk_io_time_seconds_total{job="node-exporter", instance=~"%s"}[3m]) * 100)`
)
