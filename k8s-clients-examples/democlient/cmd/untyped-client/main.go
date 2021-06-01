package main

import (
	"democlient/pkg/utils"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	guestBookGvr = schema.GroupVersionResource{
		Group:    "demo.walk1ng.io",
		Version:  "v1",
		Resource: "guestbooks",
	}

	fmtStr = "%-30s%-30s%-30s%-30s\n"
)

func main() {
	log.Println("Loading client config")
	config, err := clientcmd.BuildConfigFromFlags("", utils.UserConfig())
	utils.ExitErr("failed to load client config", err)

	log.Println("Loading dynamic client")
	dc, err := dynamic.NewForConfig(config)
	utils.ExitErr("failed to load dynamic client", err)

	log.Println("Listing GVR in default namespace")

	fmt.Printf(fmtStr, "Namespace", "Name", "Spec.Foo", "GVK")
	objs, err := dc.Resource(guestBookGvr).Namespace("default").List(v1.ListOptions{})
	if err != nil {
		utils.ExitErr("failed to list gvr", err)
	}
	for _, obj := range objs.Items {
		fmt.Printf(fmtStr,
			utils.GetResourceField(obj, "metadata", "namespace"),
			utils.GetResourceField(obj, "metadata", "name"),
			utils.GetResourceField(obj, "spec", "foo"),
			obj.GetObjectKind().GroupVersionKind().String(),
		)
	}

	log.Println("Creating guestbooks")
	err = CreateGuestBooks(dc)
	if err != nil {
		utils.ExitErr("failed to create guestbook", err)
	}

	log.Println("Creating guestbook with yaml")
	manifest := `
apiVersion: demo.walk1ng.io/v1
kind: GuestBook
metadata:
  name: gb2020
spec:
  foo: hello2020`

	err = CreateGuestBookWithManifest(dc, manifest)
	utils.ExitErr("failed to create guestbook with manifest", err)

	log.Println("Listing GVR again in default namespace")

	fmt.Printf(fmtStr, "Namespace", "Name", "Spec.Foo", "GVK")
	objs, err = dc.Resource(guestBookGvr).Namespace("default").List(v1.ListOptions{})
	if err != nil {
		utils.ExitErr("failed to list gvr", err)
	}
	for _, obj := range objs.Items {
		fmt.Printf(fmtStr,
			utils.GetResourceField(obj, "metadata", "namespace"),
			utils.GetResourceField(obj, "metadata", "name"),
			utils.GetResourceField(obj, "spec", "foo"),
			obj.GetObjectKind().GroupVersionKind().String(),
		)
	}
}

func NewGuestBook(name, foo string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "GuestBook",
			"apiVersion": guestBookGvr.Group + "/" + guestBookGvr.Version,
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": "default",
			},
			"spec": map[string]interface{}{
				"foo": foo,
			},
		},
	}
}

func CreateGuestBooks(client dynamic.Interface) error {
	guestbooks := make(map[string]string)
	guestbooks["gb2020"] = "hello2020"
	guestbooks["gb2021"] = "hello2021"

	for name, foo := range guestbooks {
		fmt.Println(name, foo)
		obj := NewGuestBook(name, foo)
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "    ")
		enc.Encode(obj)
		_, err := client.Resource(guestBookGvr).Namespace("default").Create(obj, v1.CreateOptions{})
		if err == nil || errors.IsAlreadyExists(err) {
			continue
		} else {
			return err
		}
	}

	return nil
}

func CreateGuestBookWithManifest(client dynamic.Interface, manifest string) error {
	obj := &unstructured.Unstructured{}
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, gvk, err := dec.Decode([]byte(manifest), nil, obj)
	if err != nil {
		return err
	}

	fmt.Println(gvk.GroupVersion(), gvk.Kind)

	_, err = client.Resource(guestBookGvr).Namespace("default").Create(obj, v1.CreateOptions{})

	if err == nil || errors.IsAlreadyExists(err) {
		return nil
	}
	return err
}
