package conf

const (
	QueryClusterCpuUsed  = "sum(irate(node_cpu_seconds_total{job='node-exporter', mode!='idle', instance=~'%s'}[5m]))"
	QueryClusterCpuCount = "sum(count without(cpu, mode) (node_cpu_seconds_total{job='node-exporter', mode='idle', instance=~'%s'}))"
)
