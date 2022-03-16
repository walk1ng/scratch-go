package models

import (
	corev1 "k8s.io/api/core/v1"
)

type NodeDetails struct {
	Node     *corev1.Node          `json:"node"`
	NodeInfo corev1.NodeSystemInfo `json:"nodeInfo"`
}
