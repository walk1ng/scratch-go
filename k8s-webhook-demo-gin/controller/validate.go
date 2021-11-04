package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"webhook/pkg/config"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func Validate(ar *admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
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

	// validation
	labels := deploy.ObjectMeta.Labels
	v, ok := labels[config.ReplicasValidate]
	if !ok || ok && v != "y" {
		klog.Infof("[webhook] deployment %s/%s not enable validating admission webhook\n", deploy.GetNamespace(), deploy.GetName())
		return &admissionv1.AdmissionResponse{
			Allowed: allowed,
			Result: &metav1.Status{
				Message: message,
				Code:    int32(code),
			},
		}
	}

	replicas := *deploy.Spec.Replicas
	if replicas < 3 {
		klog.Infof("[webhook] cannot create deployment due to the replicas %d less than the minimum 3\n", replicas)
		allowed = false
		code = http.StatusForbidden
		message = fmt.Sprintf("need 3 replicas at least but %d", replicas)
	}

	return &admissionv1.AdmissionResponse{
		Allowed: allowed,
		Result: &metav1.Status{
			Message: message,
			Code:    int32(code),
		},
	}
}
