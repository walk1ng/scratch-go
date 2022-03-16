package prometheus

import (
	"edge-mgr-proto/setting"
	"fmt"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

var PromV1API v1.API

func Init() error {
	client, err := api.NewClient(api.Config{
		Address: fmt.Sprintf("http://%s:%s", setting.Conf.PrometheusConfig.Host, setting.Conf.PrometheusConfig.Port),
	})
	if err != nil {
		return err
	}
	PromV1API = v1.NewAPI(client)
	return nil
}
