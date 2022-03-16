package informers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"edge-mgr-proto/setting"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	appsv1informer "k8s.io/client-go/informers/apps/v1"
	corev1informer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

var KubeInformer Informer

type Informer interface {
	Deployment() appsv1informer.DeploymentInformer
	CoreV1() corev1informer.Interface
	Run(stopper chan struct{})
}

type EdgeMgrInformer struct {
	KubeClient          kubernetes.Interface
	ResyncDuration      time.Duration
	KubeInformerFactory informers.SharedInformerFactory
}

func newEdgeMgrInformer(kubeClient kubernetes.Interface, resync time.Duration) *EdgeMgrInformer {
	inf := &EdgeMgrInformer{
		KubeClient:     kubeClient,
		ResyncDuration: resync,
	}
	inf.KubeInformerFactory = informers.NewSharedInformerFactory(kubeClient, resync)
	return inf
}

func Init(config *rest.Config) error {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// factory
	KubeInformer = newEdgeMgrInformer(client, time.Second*setting.Conf.ClientResync)

	// register deployemnt informer in sharedInformerFactory
	KubeInformer.Deployment().Informer()

	// register configmap informer in sharedInformerFactory
	KubeInformer.CoreV1().ConfigMaps().Informer()

	// reigster node informer in sharedInformerFactory
	KubeInformer.CoreV1().Nodes().Informer()

	return nil
}

func (informer *EdgeMgrInformer) Deployment() appsv1informer.DeploymentInformer {
	return informer.KubeInformerFactory.Apps().V1().Deployments()
}

func (informer *EdgeMgrInformer) CoreV1() corev1informer.Interface {
	return informer.KubeInformerFactory.Core().V1()
}

func (informer *EdgeMgrInformer) Run(stopper chan struct{}) {
	// run deployment informer
	// zap.L().Info("starting deployment informer")
	// go informer.Deployment().Informer().Run(stopper)
	// run configmap informer
	// zap.L().Info("starting configmap informer")
	// go informer.CoreV1().ConfigMaps().Informer().Run(stopper)

	// start the factory
	informer.KubeInformerFactory.Start(stopper)
	// wait for index cache synced
	informer.KubeInformerFactory.WaitForCacheSync(stopper)
}

func (informer *EdgeMgrInformer) RegisterDeploymentEventWatcher() {

	deploymentLister := informer.Deployment().Lister()

	informer.Deployment().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// method 1, assert obj as metav1.Object
			// and use its method GetNamespace() and GetName()
			deploy := obj.(metav1.Object)
			zap.L().Info("new deployment added", zap.String("namespace", deploy.GetNamespace()), zap.String("name", deploy.GetName()))
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {

			// assert object to corev1.Pod
			oldDeploy := oldObj.(*appsv1.Deployment)
			newDeploy := newObj.(*appsv1.Deployment)

			// get the desired pod to use pod lister to fetch from local indexed store
			name, namespace := oldDeploy.GetName(), oldDeploy.GetNamespace()
			oldDeploy1, err := deploymentLister.Deployments(namespace).Get(name)
			if err != nil {
				utilruntime.HandleError(errors.WithMessagef(err, "deployment lister to get deployment %s/%s failed", namespace, name))
				return
			}

			if oldDeploy.GetResourceVersion() == newDeploy.GetResourceVersion() {
				zap.L().Info("old deployment updated due to caches synced", zap.String("namespace", oldDeploy.GetNamespace()), zap.String("name", oldDeploy.GetName()))

				prettyPrint("LOCAL STORE", oldDeploy1)
				return
			}

			zap.L().Info("old deployment updated", zap.String("namespace", oldDeploy.GetNamespace()), zap.String("name", oldDeploy.GetName()))

			prettyPrint("BEFORE", oldDeploy)
			prettyPrint("AFTER", newDeploy)
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err != nil {
				utilruntime.HandleError(err)
				return
			}

			// method 2: use cache.SplitMetaNamespaceKey to get namespace and name of object
			namespace, name, err := cache.SplitMetaNamespaceKey(key)
			if err != nil {
				utilruntime.HandleError(errors.WithMessage(err, fmt.Sprintf("invalid resource key: %s", key)))
				return
			}

			zap.L().Info("old deployment deleted", zap.String("namespace", namespace), zap.String("name", name))
		},
	})
}

func prettyPrint(title string, obj interface{}) {
	log.Printf(" =============  %s =============\n", title)
	data, _ := json.MarshalIndent(obj, "", "    ")
	log.Println(string(data))
	log.Printf(" =============  %s =============\n", title)
}
