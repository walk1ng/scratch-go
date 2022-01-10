package k8sunittestdemo

import (
	"context"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

func TestAPI_NewPodWithMeta(t *testing.T) {
	type fields struct {
		Client kubernetes.Interface
	}
	type args struct {
		namespace string
		name      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "test pod1",
			fields:  fields{Client: fake.NewSimpleClientset()},
			args:    args{namespace: "ns1", name: "pod1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := API{
				Client: tt.fields.Client,
			}
			err := api.NewPodWithMeta(tt.args.namespace, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.NewPodWithMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// get the created pod
			if got, err := api.Client.CoreV1().Pods(tt.args.namespace).Get(context.Background(), tt.args.name, metav1.GetOptions{}); err != nil {
				t.Errorf("failed to get created pod, error: %v\n", err)
			} else if got.Name != tt.args.name {
				t.Errorf("pod name err: expected: %s but got %s\n", tt.args.name, got.Name)
			} else if got.Namespace != tt.args.namespace {
				t.Errorf("pod namespace err: expected: %s but got %s\n", tt.args.namespace, got.Namespace)
			}
		})
	}
}

func TestCache_AddPodInCache(t *testing.T) {
	type fields struct {
		Client kubernetes.Interface
		Pods   map[string]*corev1.Pod
	}
	type args struct {
		p *corev1.Pod
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test add pod in cache",
			fields: fields{
				Client: fake.NewSimpleClientset(),
				Pods:   make(map[string]*corev1.Pod),
			},
			args: args{
				p: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name: "testpod",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			c := Cache{
				Client: tt.fields.Client,
				Pods:   tt.fields.Pods,
			}

			// a channel to hold the added pod
			pods := make(chan *corev1.Pod)

			// pod informers
			informers := informers.NewSharedInformerFactory(c.Client, time.Millisecond*1000)
			podInformer := informers.Core().V1().Pods().Informer()
			podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					pod := obj.(*corev1.Pod)
					// in channel
					pods <- pod
					t.Logf("pod added in channel: %s/%s", pod.Namespace, pod.Name)
					// in cache
					if err := c.AddPodInCache(pod); nil != err {
						t.Errorf("add pod in cache err %v", err)
					}
				}})

			err := c.AddPodInCache(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cache.AddPodInCache() error = %v, wantErr %v", err, tt.wantErr)
			}

			// start informer
			informers.Start(ctx.Done())

			// sync and list pods before test
			cache.WaitForCacheSync(ctx.Done(), podInformer.HasSynced)

			// create a test pod
			_, err = c.Client.CoreV1().Pods("default").Create(context.Background(), tt.args.p, metav1.CreateOptions{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Cache.AddPodInCache() error = %v, wantErr %v", err, tt.wantErr)
			}

			select {
			case pod := <-pods:
				if _, ok := c.Pods[pod.Name]; !ok {
					t.Errorf("No pod added in cache after create pod in k8s\n")
				}
			case <-time.After(wait.ForeverTestTimeout):
				t.Errorf("pod informer can't notify pod creation")
			}

		})
	}
}
