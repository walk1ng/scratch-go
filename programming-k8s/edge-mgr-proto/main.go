package main

import (
	"edge-mgr-proto/mock"
	"edge-mgr-proto/mq"
	"edge-mgr-proto/routers"
	"edge-mgr-proto/service"
	"edge-mgr-proto/setting"
	"fmt"
	"os"

	"edge-mgr-proto/common"
	"edge-mgr-proto/logger"

	"edge-mgr-proto/pkg/client"
	"edge-mgr-proto/pkg/informers"

	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("config file is required, e.g.: ./edge-mgr-proto config.yaml")
		return
	}

	// init setting
	if err := setting.Init(os.Args[1]); err != nil {
		fmt.Printf("failed to load config file: %v\n", err)
		return
	}

	// init logger
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Printf("failed to init logger: %v\n", err)
		return
	}

	// init routers
	r := routers.SetupRouter(setting.Conf.Mode)

	// rest config
	var config *rest.Config
	var err error
	config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	}

	// init kube clients
	if err := client.Init(config); err != nil {
		zap.L().Panic("failed to init kube clients", zap.Error(err))
	}

	// init informers
	if err := informers.Init(config); err != nil {
		zap.L().Panic("failed to init kube informers", zap.Error(err))
	}

	// register watcher for deployments
	// informers.KubeInformer.RegisterDeploymentEventWatcher()
	// run informers
	stopper := make(chan struct{})
	informers.KubeInformer.Run(stopper)

	// init channels
	common.WorkChan = &mq.WorkChannel{
		Queue: make(chan mq.Message, 10),
	}

	// init service
	if err := service.Init(common.WorkChan, informers.KubeInformer, client.KubeClient); err != nil {
		zap.L().Panic("failed to init service", zap.Error(err))
	}

	// init mock cluster db
	if err := mock.InitClusterDB(); err != nil {
		zap.L().Panic("failed to init mock cluster db", zap.Error(err))
	}

	// run server
	if err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port)); err != nil {
		zap.L().Fatal("failed to run server", zap.Error(err))
		return
	}
}
