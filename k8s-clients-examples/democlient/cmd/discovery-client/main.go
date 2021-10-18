package main

import (
	"democlient/pkg/utils"
	"fmt"
	"log"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.Println("Loading client config")
	config, err := clientcmd.BuildConfigFromFlags("", utils.UserConfig())
	utils.ExitErr("failed to load client config", err)

	log.Println("Loading discovery client")
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	utils.ExitErr("failed to load discovery client", err)

	_, apiResourceList, err := dc.ServerGroupsAndResources()
	// dc.ServerResources()
	if err != nil {
		panic(err)
	}

	for _, list := range apiResourceList {
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			panic(err)
		}

		for _, resource := range list.APIResources {
			fmt.Printf("name: %s, group: %s, version: %s\n", resource.Name, gv.Group, gv.Version)
		}
	}
}
