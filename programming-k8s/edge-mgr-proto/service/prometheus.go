package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"edge-mgr-proto/conf"
	"edge-mgr-proto/types"

	"github.com/pkg/errors"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"go.uber.org/zap"
)

type PrometheusService interface {
	GetClusterCpuUsage(nodeList []string) *types.ClusterCpuUsage
	// GetClusterCpuUsageRange(nodeList []string)
	GetClusterMemoryUsage(nodeList []string) *types.ClusterMemUsage
	// GetClusterMemoryUsageRange(nodeList []string)
	GetClusterDiskUsage(nodeList []string) *types.ClusterDiskUsage
	// GetClusterDiskUsageRange(nodeList []string)
	GetNodeCpuUsage(node string) (model.SampleValue, error)
	GetNodeMemoryUsage(node string) (model.SampleValue, error)
	GetNodeDiskUsage(node string) (model.SampleValue, error)
	GetNodeDiskIOUsage(node string) (model.SampleValue, error)
}

type prometheusService struct {
	v1.API
}

func newPrometheusService(api v1.API) *prometheusService {
	return &prometheusService{
		api,
	}
}

/*
  Cluster metrics
*/

func (svc *prometheusService) GetClusterCpuUsage(nodeList []string) *types.ClusterCpuUsage {
	instances := strings.Join(nodeList, "|")
	clusterCpuUsedQuery := fmt.Sprintf(conf.QueryClusterCpuUsed, instances)
	clusterCpuCountQuery := fmt.Sprintf(conf.QueryClusterCpuCount, instances)

	// used
	zap.L().Info("getting metrics usage", zap.String("nodes", instances), zap.Any("target", types.TargetClusterUsedCPU))
	used, err := svc.query(context.Background(), clusterCpuUsedQuery)
	if err != nil {
		zap.L().Error("error getting metrics usage", zap.Error(err), zap.String("nodes", instances), zap.Any("target", types.TargetClusterUsedCPU))
	}

	// count
	zap.L().Info("getting metrics usage", zap.String("nodes", instances), zap.Any("target", types.TargetClusterTotalCPU))
	count, err := svc.query(context.Background(), clusterCpuCountQuery)
	if err != nil {
		zap.L().Error("error getting metrics usage", zap.Error(err), zap.String("nodes", instances), zap.Any("target", types.TargetClusterTotalCPU))
	}

	return &types.ClusterCpuUsage{
		Used:  float64(getFirstValue(used)),
		Total: float64(getFirstValue(count)),
	}
}

/* func (svc *prometheusService) GetClusterCpuUsageRange(nodeList []string) {
	panic("not implement")
} */

func (svc *prometheusService) GetClusterMemoryUsage(nodeList []string) *types.ClusterMemUsage {
	instances := strings.Join(nodeList, "|")
	clusterMemUsedQuery := fmt.Sprintf(conf.QueryClusterMemoryUsed, instances, instances, instances, instances, instances)
	clusterMemTotalQuery := fmt.Sprintf(conf.QueryClusterMemoryTotal, instances)

	// used
	zap.L().Info("getting metrics usage", zap.String("nodes", instances), zap.Any("target", types.TargetClusterUsedMemory))
	used, err := svc.query(context.Background(), clusterMemUsedQuery)
	if err != nil {
		zap.L().Error("error getting metrics usage", zap.Error(err), zap.String("nodes", instances), zap.Any("target", types.TargetClusterUsedMemory))
	}

	// total
	zap.L().Info("getting metrics usage", zap.String("nodes", instances), zap.Any("target", types.TargetClusterTotalMemory))
	total, err := svc.query(context.Background(), clusterMemTotalQuery)
	if err != nil {
		zap.L().Error("error getting metrics usage", zap.Error(err), zap.String("nodes", instances), zap.Any("target", types.TargetClusterTotalMemory))
	}

	return &types.ClusterMemUsage{
		BytesUsed:  float64(getFirstValue(used)),
		BytesTotal: float64(getFirstValue(total)),
	}
}

/* func (svc *prometheusService) GetClusterMemoryUsageRange(nodeList []string) {
	panic("not implement")
} */

func (svc *prometheusService) GetClusterDiskUsage(nodeList []string) *types.ClusterDiskUsage {
	instances := strings.Join(nodeList, "|")
	clusterDiskUsedQuery := fmt.Sprintf(conf.QueryClusterDiskUsed, instances, instances)
	clusterDiskTotalQuery := fmt.Sprintf(conf.QueryClusterDiskTotal, instances)

	// used
	zap.L().Info("getting metrics usage", zap.String("nodes", instances), zap.Any("target", types.TargetClusterUsedDisk))
	used, err := svc.query(context.Background(), clusterDiskUsedQuery)
	if err != nil {
		zap.L().Error("error getting metrics usage", zap.Error(err), zap.String("nodes", instances), zap.Any("target", types.TargetClusterUsedDisk))
	}

	// total
	zap.L().Info("getting metrics usage", zap.String("nodes", instances), zap.Any("target", types.TargetClusterTotalDisk))
	total, err := svc.query(context.Background(), clusterDiskTotalQuery)
	if err != nil {
		zap.L().Error("error getting metrics usage", zap.Error(err), zap.String("nodes", instances), zap.Any("target", types.TargetClusterTotalDisk))
	}

	return &types.ClusterDiskUsage{
		BytesUsed:  float64(getFirstValue(used)),
		BytesTotal: float64(getFirstValue(total)),
	}
}

/* func (svc *prometheusService) GetClusterDiskUsageRange(nodeList []string) {
	panic("not implement")
} */

/*
  Node metrics
*/

func (svc *prometheusService) GetNodeCpuUsage(node string) (model.SampleValue, error) {
	zap.L().Info("getting metrics usage", zap.String("node", node), zap.Any("target", types.TargetNodeCPU))
	nodeCpuUsageQuery := fmt.Sprintf(conf.QueryNodeCpuUsage, node, node)
	usage, err := svc.query(context.Background(), nodeCpuUsageQuery)
	if err != nil {
		zap.L().Error("error getting metrics usage", zap.Error(err), zap.String("node", node), zap.Any("target", types.TargetNodeCPU))
		return 0, errors.WithMessagef(err, "get %s usage of node %s", types.TargetNodeCPU, node)
	}
	return getFirstValue(usage), nil
}

func (svc *prometheusService) GetNodeMemoryUsage(node string) (model.SampleValue, error) {
	zap.L().Info("getting metrics usage", zap.String("node", node), zap.Any("target", types.TargetNodeMemory))
	nodeMemUsageQuery := fmt.Sprintf(conf.QueryNodeMemoryUsage, node, node, node, node, node, node)
	usage, err := svc.query(context.Background(), nodeMemUsageQuery)
	if err != nil {
		zap.L().Error("error getting metrics usage", zap.Error(err), zap.String("node", node), zap.Any("target", types.TargetNodeMemory))
		return 0, errors.WithMessagef(err, "get %s usage of node %s", types.TargetNodeMemory, node)
	}
	return getFirstValue(usage), nil
}

func (svc *prometheusService) GetNodeDiskUsage(node string) (model.SampleValue, error) {
	zap.L().Info("getting metrics usage", zap.String("node", node), zap.Any("target", types.TargetNodeDisk))
	nodeDiskUsageQuery := fmt.Sprintf(conf.QueryNodeDiskUsage, node, node)
	usage, err := svc.query(context.Background(), nodeDiskUsageQuery)
	if err != nil {
		zap.L().Error("error getting metrics usage", zap.Error(err), zap.String("node", node), zap.Any("target", types.TargetNodeDisk))
		return 0, errors.WithMessagef(err, "get %s usage of node %s", types.TargetNodeDisk, node)
	}
	return getFirstValue(usage), nil
}

func (svc *prometheusService) GetNodeDiskIOUsage(node string) (model.SampleValue, error) {
	zap.L().Info("getting metrics usage", zap.String("node", node), zap.Any("target", types.TargetNodeDiskIO))
	nodeDiskIOUsageQuery := fmt.Sprintf(conf.QueryNodeDiskIOUsage, node)
	usage, err := svc.query(context.Background(), nodeDiskIOUsageQuery)
	if err != nil {
		zap.L().Error("error getting metrics usage", zap.Error(err), zap.String("node", node), zap.Any("target", types.TargetNodeDiskIO))
		return 0, errors.WithMessagef(err, "get %s usage of node %s", types.TargetNodeDiskIO, node)
	}
	return getFirstValue(usage), nil
}

func (svc *prometheusService) query(c context.Context, q string) (model.Value, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	result, warnings, err := svc.Query(ctx, q, time.Now())
	if err != nil {
		zap.L().Error("error querying prometheus", zap.String("query", q), zap.Error(err))
		return nil, err
	}
	if len(warnings) > 0 {
		zap.L().Warn("warning querying prometheus", zap.String("query", q), zap.Int("warn", len(warnings)))
	}

	zap.L().Info("querying prometheus", zap.String("query", q), zap.Any("result", result))
	return result, nil
}

func (svc *prometheusService) queryRange(c context.Context, q, start, end, step string) {
	panic("not implement")
}

func getFirstValue(in model.Value) model.SampleValue {
	// vector
	iin := in.(model.Vector)
	zap.L().Info("getFirstValue", zap.Int("len(model.Value)", iin.Len()))
	if iin.Len() == 0 {
		zap.L().Info("getFirstValue: return zero value")
		return 0
	}
	return iin[0].Value
}
