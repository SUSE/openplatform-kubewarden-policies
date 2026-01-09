package domain

import (
	apimachinery_pkg_apis_meta_v1 "github.com/kubewarden/k8s-objects/apimachinery/pkg/apis/meta/v1"
)

type Metadata struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type PodSpecResources struct {
	Limits   map[string]interface{} `json:"limits"`
	Requests map[string]interface{} `json:"requests"`
}

type ContainerSpec struct {
	Resources PodSpecResources `json:"resources"`
}

type PodSpec struct {
	Containers []ContainerSpec `json:"containers"`
}

type Pod struct {
	Metadata Metadata `json:"metadata"`
	Spec     PodSpec  `json:"spec"`
}

type ResourceQuotaSpec struct {
	Hard map[string]interface{} `json:"hard"`
}

type ResourceQuota struct {
	Spec ResourceQuotaSpec `json:"spec"`
}

type ResourceQuotaList struct {
	APIVersion string                                  `json:"apiVersion,omitempty"`
	Items      []ResourceQuota                         `json:"items"`
	Kind       string                                  `json:"kind,omitempty"`
	Metadata   *apimachinery_pkg_apis_meta_v1.ListMeta `json:"metadata,omitempty"`
}
