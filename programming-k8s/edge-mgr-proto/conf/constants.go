package conf

const (
	// CPU
	QueryClusterCpuUsed  = `sum(irate(node_cpu_seconds_total{job="node-exporter", mode!="idle", instance=~"%s"}[5m]))`
	QueryClusterCpuCount = `sum(count without(cpu, mode) (node_cpu_seconds_total{job="node-exporter", mode="idle", instance=~"%s"}))`

	// MEMORY
	QueryClusterMemoryTotal = `sum(node_memory_MemTotal_bytes{job="node-exporter", instance=~"%s"})`
	QueryClusterMemoryUsed  = `sum(node_memory_MemTotal_bytes{job="node-exporter", instance=~"%s"}) - 
	sum(node_memory_MemFree_bytes{job="node-exporter", instance=~"%s"}) - 
	sum(node_memory_Buffers_bytes{job="node-exporter", instance=~"%s"}) - 
	sum(node_memory_Cached_bytes{job="node-exporter", instance=~"%s"}) + 
	sum(node_memory_Shmem_bytes{job="node-exporter", instance=~"%s"})`

	// DISK
	// TODO
	QueryClusterDiskTotal = `sum(node_filesystem_size_bytes{job="node-exporter", instance=~"%s", fstype=~"{ DISK_FSTYPE }", mountpoint=~"{ DISK_MOUNTPOINT }"}})`
	QueryClusterDiskUsed  = `sum(node_filesystem_size_bytes{{cluster_id="{cluster_id}", job="node-exporter", instance=~"{node_ip_list}", fstype=~"{ DISK_FSTYPE }", mountpoint=~"{ DISK_MOUNTPOINT }"}}) - sum(node_filesystem_free_bytes{{cluster_id="{cluster_id}", job="node-exporter", instance=~"{node_ip_list}", fstype=~"{ DISK_FSTYPE }", mountpoint=~"{ DISK_MOUNTPOINT }"}})`
)
