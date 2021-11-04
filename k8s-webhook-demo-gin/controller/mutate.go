package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"webhook/pkg/config"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/klog/v2"
)

type patchOption struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

type sideConfig struct {
	Containers []corev1.Container `yaml:"containers"`
}

func loadSideCarConfig() (*sideConfig, error) {
	data, err := ioutil.ReadFile(config.SideCarConfigFile)
	if err != nil {
		klog.Errorf("[webhook] cannot load sidecar configuration: %v\n", err)
		return nil, err
	}

	sideCar := sideConfig{}
	err = yaml.Unmarshal(data, &sideCar)
	if err != nil {
		klog.Errorf("[webhook] cannot parse sidecar configuration: %v\n", err)
		return nil, err
	}

	return &sideCar, nil
}

func deploymentPatch(deploy *appsv1.Deployment, sideCar *sideConfig) ([]byte, error) {
	var patch []patchOption
	patch = append(patch, addContainers(deploy.Spec.Template.Spec.Containers, sideCar.Containers, "/spec/template/spec/containers")...)
	return json.Marshal(patch)
}

func addContainers(target, added []corev1.Container, basePath string) (patch []patchOption) {
	first := len(target) == 0
	var value interface{}
	for _, add := range added {
		value = add
		path := basePath
		if first {
			first = false
			value = []corev1.Container{add}
		} else {
			path = basePath + "/-"
		}

		patch = append(patch, patchOption{
			Op:    "add",
			Path:  path,
			Value: value,
		})
	}

	return
}

func Mutate(ar *admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	req := ar.Request
	var (
		allowed bool   = true
		code    int    = http.StatusOK
		message string = ""
	)

	// unmarshal the raw request object to a deployment
	var deploy appsv1.Deployment
	err := json.Unmarshal(req.Object.Raw, &deploy)
	if err != nil {
		klog.Errorf("[webhook] cannot unmarshal AdmissionReview Request object raw: %v\n", err)
		allowed = false
		code = http.StatusBadRequest
		return &admissionv1.AdmissionResponse{
			Allowed: allowed,
			Result: &metav1.Status{
				Message: err.Error(),
				Code:    int32(code),
			},
		}
	}

	labels := deploy.ObjectMeta.Labels
	v, ok := labels[config.SideCarInject]
	if !ok || ok && v != "y" {
		klog.Infof("[webhook] deployment %s/%s not enable mutating admission webhook\n", deploy.GetNamespace(), deploy.GetName())
		return &admissionv1.AdmissionResponse{
			Allowed: allowed,
			Result: &metav1.Status{
				Message: message,
				Code:    int32(code),
			},
		}
	}

	// load the sidecar configuration
	sideCar, err := loadSideCarConfig()
	if err != nil {
		allowed = false
		code = http.StatusBadRequest
		return &admissionv1.AdmissionResponse{
			Allowed: allowed,
			Result: &metav1.Status{
				Message: err.Error(),
				Code:    int32(code),
			},
		}
	}

	// generate deployment patch
	patchBytes, err := deploymentPatch(&deploy, sideCar)
	if err != nil {
		allowed = false
		code = http.StatusBadRequest
		return &admissionv1.AdmissionResponse{
			Allowed: allowed,
			Result: &metav1.Status{
				Message: err.Error(),
				Code:    int32(code),
			},
		}
	}

	klog.Infof("[webhook] AdmissionResponse patch=%v\n", string(patchBytes))
	return &admissionv1.AdmissionResponse{
		Allowed: allowed,
		Patch:   patchBytes,
		PatchType: func() *admissionv1.PatchType {
			patchType := admissionv1.PatchTypeJSONPatch
			return &patchType
		}(),
		Result: &metav1.Status{
			Code:    int32(code),
			Message: message,
		},
	}
}
