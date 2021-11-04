package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"webhook/controller"

	"github.com/gin-gonic/gin"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/klog/v2"
)

var (
	runtimeScheme = runtime.NewScheme()
	codeFactory   = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codeFactory.UniversalDeserializer()
)

func ValidatingAdmission(c *gin.Context) {
	// parse request
	reqData, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		klog.Errorf("[webhook] failed to parse request body: %v\n", err)
		http.Error(c.Writer, fmt.Sprintf("[webhook] failed to parse request body: %v", err), http.StatusBadRequest)
		return
	}

	// parse AdmissionReview
	var requestedAdmissionReview admissionv1.AdmissionReview
	var admissionResponse *admissionv1.AdmissionResponse
	_, _, err = deserializer.Decode(reqData, nil, &requestedAdmissionReview)
	if err != nil {
		klog.Errorf("[webhook] failed to parse request body: %v\n", err)
		admissionResponse = &admissionv1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		}
	} else {
		// validate AdmissionReview and get AdmissionReview response
		admissionResponse = controller.Validate(&requestedAdmissionReview)
	}

	// construct the responsed AdmissionReview, it will be sent to apiserver again
	responsedAdmissionReview := admissionv1.AdmissionReview{}
	responsedAdmissionReview.APIVersion = requestedAdmissionReview.APIVersion
	responsedAdmissionReview.Kind = requestedAdmissionReview.Kind

	if admissionResponse != nil {
		responsedAdmissionReview.Response = admissionResponse
		if requestedAdmissionReview.Request != nil {
			responsedAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID
		}
	}

	klog.Infof("[webhook] apiserver callback, send %v\n", responsedAdmissionReview)

	respData, err := json.Marshal(responsedAdmissionReview)
	if err != nil {
		klog.Errorf("[webhook] cannot parse apiserver callback AdmissionReview: %v\n", err)
		http.Error(c.Writer, fmt.Sprintf("[webhook] cannot parse apiserver callback AdmissionReview: %v", err), http.StatusInternalServerError)
		return
	}

	klog.Infoln("[webhook] ready to send apiserver callback AdmissionReview")
	if _, err := c.Writer.Write(respData); err != nil {
		klog.Errorf("[webhook] failed to send apiserver callback AdmissionReview: %v\n", err)
		http.Error(c.Writer, fmt.Sprintf("[webhook] failed to send apiserver callback AdmissionReview: %v", err), http.StatusInternalServerError)
		return
	}

	klog.Infoln("[webhook] sent apiserver callback AdmissionReview done")
}

func MutatingAdmission(c *gin.Context) {
	// parse request
	reqData, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		klog.Errorf("[webhook] failed to parse request body: %v\n", err)
		http.Error(c.Writer, fmt.Sprintf("[webhook] failed to parse request body: %v", err), http.StatusBadRequest)
		return
	}

	// parse AdmissionReview
	var requestedAdmissionReview admissionv1.AdmissionReview
	var admissionResponse *admissionv1.AdmissionResponse
	_, _, err = deserializer.Decode(reqData, nil, &requestedAdmissionReview)
	if err != nil {
		klog.Errorf("[webhook] failed to parse request body: %v\n", err)
		admissionResponse = &admissionv1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		}
	} else {
		// validate AdmissionReview and get AdmissionReview response
		admissionResponse = controller.Mutate(&requestedAdmissionReview)
	}

	// construct the responsed AdmissionReview, it will be sent to apiserver again
	responsedAdmissionReview := &admissionv1.AdmissionReview{}
	responsedAdmissionReview.APIVersion = requestedAdmissionReview.APIVersion
	responsedAdmissionReview.Kind = requestedAdmissionReview.Kind
	if admissionResponse != nil {
		responsedAdmissionReview.Response = admissionResponse
		if requestedAdmissionReview.Request != nil {
			responsedAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID
		}
	}

	klog.Infof("[webhook] apiserver callback, send %v\n", responsedAdmissionReview)

	respData, err := json.Marshal(responsedAdmissionReview)
	if err != nil {
		klog.Errorf("[webhook] cannot parse apiserver callback AdmissionReview: %v\n", err)
		http.Error(c.Writer, fmt.Sprintf("[webhook] cannot parse apiserver callback AdmissionReview: %v", err), http.StatusInternalServerError)
		return
	}

	klog.Infoln("[webhook] ready to send apiserver callback AdmissionReview")
	if _, err := c.Writer.Write(respData); err != nil {
		klog.Errorf("[webhook] failed to send apiserver callback AdmissionReview: %v\n", err)
		http.Error(c.Writer, fmt.Sprintf("[webhook] failed to send apiserver callback AdmissionReview: %v", err), http.StatusInternalServerError)
		return
	}

	klog.Infoln("[webhook] sent apiserver callback AdmissionReview done")
}
