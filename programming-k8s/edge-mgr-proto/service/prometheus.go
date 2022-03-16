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
	GetClusterCpuUsage(nodeList []string) map[string]interface{}
	GetClusterCpuUsageRange(nodeIPList []string)
	GetClusterMemoryUsage(nodeIPList []string)
	GetClusterMemoryUsageRange(nodeIPList []string)
	GetClusterDiskUsage(nodeIPList []string)
	GetClusterDiskUsageRange(nodeIPList []string)
}

type prometheusService struct {
	v1.API
}

func newPrometheusService(api v1.API) *prometheusService {
	return &prometheusService{
		api,
	}
}

func (svc *prometheusService) GetClusterCpuUsage(nodeList []string) map[string]interface{} {
	instances := strings.Join(nodeList, "|")
	clusterCpuUsedQuery := fmt.Sprintf(conf.QueryClusterCpuUsed, instances)
	clusterCpuCountQuery := fmt.Sprintf(conf.QueryClusterCpuCount, instances)

	used := svc.query(context.Background(), clusterCpuUsedQuery).(model.Vector)
	count := svc.query(context.Background(), clusterCpuCountQuery).(model.Vector)

	fmt.Printf("used: %+v\n", used[0].Value.String())

	return map[string]interface{}{
		"used":  used[0].Value.String(),
		"total": count[0].Value.String(),
	}
}

func (svc *prometheusService) GetClusterCpuUsageRange(nodeIPList []string) {

}

func (svc *prometheusService) GetClusterMemoryUsage(nodeIPList []string) {

}

func (svc *prometheusService) GetClusterMemoryUsageRange(nodeIPList []string) {

}

func (svc *prometheusService) GetClusterDiskUsage(nodeIPList []string) {

}

func (svc *prometheusService) GetClusterDiskUsageRange(nodeIPList []string) {

}

func (svc *prometheusService) query(c context.Context, q string) model.Value {
	fmt.Println(q)
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	// used
	result, warnings, err := svc.Query(ctx, q, time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	fmt.Printf("Result:\n%v\n", result)

	return result
}

func (svc *prometheusService) queryRange(c context.Context, q, start, end, step string) {

}
