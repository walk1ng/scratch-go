package utils

import (
	"log"
	"os/user"
	"path/filepath"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func UserConfig() string {
	usr, err := user.Current()
	if err != nil {
		ExitErr("failed to get user", err)
	}
	return filepath.Join([]string{usr.HomeDir, ".kube", "config"}...)
}

func ExitErr(msg string, err error) {
	if err != nil {
		log.Fatalf("%s:%v\n", msg, err)
	}
}

func GetResourceField(obj unstructured.Unstructured, fieldPath ...string) string {
	fieldValue, found, err := unstructured.NestedString(obj.Object, fieldPath...)
	if err != nil || !found {
		return ""
	}

	return fieldValue
}
