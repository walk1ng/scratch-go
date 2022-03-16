package prometheus

import (
	"context"
	"edge-mgr-proto/setting"
	"fmt"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

var PromClient api.Client

func Init() error {
	var err error
	PromClient, err = api.NewClient(api.Config{
		// Address: "http://demo.robustperception.io:9090",
		Address: fmt.Sprintf("http://%s:%s", setting.Conf.PrometheusConfig.Host, setting.Conf.PrometheusConfig.Port),
	})
	if err != nil {
		return err
	}
	return nil
}

func GetClusterCpuUsage(nodeIPList []string) {
	v1api := v1.NewAPI(PromClient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cpuUsedQuery := ""
	cpuCountQuery := ""

	// used
	result, warnings, err := v1api.Query(ctx, cpuUsedQuery, time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	fmt.Printf("Result:\n%v\n", result)

	// used
	result1, warnings1, err1 := v1api.Query(ctx, cpuCountQuery, time.Now())
	if err1 != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err1)
		os.Exit(1)
	}
	if len(warnings1) > 0 {
		fmt.Printf("Warnings: %v\n", warnings1)
	}
	fmt.Printf("Result:\n%v\n", result1)
}

func GetClusterCpuUsageRange(nodeIPList []string) {

}

func GetClusterMemoryUsage(nodeIPList []string) {

}

func GetClusterMemoryUsageRange(nodeIPList []string) {

}

func GetClusterDiskUsage(nodeIPList []string) {

}

func GetClusterDiskUsageRange(nodeIPList []string) {

}

func query(q string) {

}

func queryRange(q, start, end, step string) {

}
