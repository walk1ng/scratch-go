package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"edge-mgr-proto/conf"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type PrometheusService interface {
	GetClusterCpuUsage(nodeList []string) map[string]model.SampleValue
	// GetClusterCpuUsageRange(nodeList []string)
	GetClusterMemoryUsage(nodeList []string) map[string]model.SampleValue
	// GetClusterMemoryUsageRange(nodeList []string)
	GetClusterDiskUsage(nodeList []string) map[string]model.SampleValue
	// GetClusterDiskUsageRange(nodeList []string)
	GetNodeCpuUsage(node string) model.SampleValue
	GetNodeMemoryUsage(node string) model.SampleValue
	GetNodeDiskUsage(node string) model.SampleValue
	GetNodeDiskIOUsage(node string) model.SampleValue
}

type prometheusService struct {
	v1.API
}

func newPrometheusService(api v1.API) *prometheusService {
	return &prometheusService{
		api,
	}
}

func (svc *prometheusService) GetClusterCpuUsage(nodeList []string) map[string]model.SampleValue {
	instances := strings.Join(nodeList, "|")
	clusterCpuUsedQuery := fmt.Sprintf(conf.QueryClusterCpuUsed, instances)
	clusterCpuCountQuery := fmt.Sprintf(conf.QueryClusterCpuCount, instances)

	used := svc.query(context.Background(), clusterCpuUsedQuery).(model.Vector)
	count := svc.query(context.Background(), clusterCpuCountQuery).(model.Vector)

	fmt.Printf("used: %+v\n", used[0].Value.String())

	return map[string]model.SampleValue{
		"used":  used[0].Value,
		"total": count[0].Value,
	}

}

/* func (svc *prometheusService) GetClusterCpuUsageRange(nodeList []string) {
	panic("not implement")
} */

func (svc *prometheusService) GetClusterMemoryUsage(nodeList []string) map[string]model.SampleValue {
	instances := strings.Join(nodeList, "|")
	clusterMemUsedQuery := fmt.Sprintf(conf.QueryClusterMemoryUsed, instances, instances, instances, instances, instances)
	clusterMemTotalQuery := fmt.Sprintf(conf.QueryClusterMemoryTotal, instances)

	used := svc.query(context.Background(), clusterMemUsedQuery).(model.Vector)
	total := svc.query(context.Background(), clusterMemTotalQuery).(model.Vector)

	fmt.Printf("used: %+v\n", used[0].Value.String())

	return map[string]model.SampleValue{
		"used_bytes":  used[0].Value,
		"total_bytes": total[0].Value,
	}

}

/* func (svc *prometheusService) GetClusterMemoryUsageRange(nodeList []string) {
	panic("not implement")
} */

func (svc *prometheusService) GetClusterDiskUsage(nodeList []string) map[string]model.SampleValue {
	// TODO
	/* instances := strings.Join(nodeList, "|")
	clusterDiskUsedQuery := fmt.Sprintf(conf.QueryClusterDiskUsed, instances, instances, instances, instances, instances)
	clusterDiskTotalQuery := fmt.Sprintf(conf.QueryClusterDiskTotal, instances)

	used := svc.query(context.Background(), clusterDiskUsedQuery).(model.Vector)
	total := svc.query(context.Background(), clusterDiskTotalQuery).(model.Vector)

	fmt.Printf("used: %+v\n", used[0].Value.String())

	return map[string]model.SampleValue{
		"used_bytes":  used[0].Value,
		"total_bytes": total[0].Value,
	} */
	panic("not implement")
}

/* func (svc *prometheusService) GetClusterDiskUsageRange(nodeList []string) {
	panic("not implement")
} */

func (svc *prometheusService) GetNodeCpuUsage(node string) model.SampleValue {
	nodeCpuUsageQuery := fmt.Sprintf(conf.QueryNodeCpuUsage, node, node)
	usage := svc.query(context.Background(), nodeCpuUsageQuery)
	return getFirstValue(usage)
}

func (svc *prometheusService) GetNodeMemoryUsage(node string) model.SampleValue {
	nodeMemUsageQuery := fmt.Sprintf(conf.QueryNodeMemoryUsage, node, node, node, node, node, node)
	usage := svc.query(context.Background(), nodeMemUsageQuery)
	return getFirstValue(usage)
}

func (svc *prometheusService) GetNodeDiskUsage(node string) model.SampleValue {
	nodeDiskUsageQuery := fmt.Sprintf(conf.QueryNodeDiskUsage, node, node)
	usage := svc.query(context.Background(), nodeDiskUsageQuery)
	return getFirstValue(usage)
}

func (svc *prometheusService) GetNodeDiskIOUsage(node string) model.SampleValue {
	nodeDiskIOUsageQuery := fmt.Sprintf(conf.QueryNodeDiskIOUsage, node)
	usage := svc.query(context.Background(), nodeDiskIOUsageQuery)
	return getFirstValue(usage)
}

func (svc *prometheusService) query(c context.Context, q string) model.Value {
	fmt.Println(q)
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	result, warnings, err := svc.Query(ctx, q, time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	fmt.Printf("Result:\n%+v\n", result)

	return result
}

func (svc *prometheusService) queryRange(c context.Context, q, start, end, step string) {
	panic("not implement")
}

func getFirstValue(in model.Value) model.SampleValue {
	// vector
	iin := in.(model.Vector)
	if iin.Len() == 0 {
		println("oh no")
		return 0
	}
	return iin[0].Value
}
