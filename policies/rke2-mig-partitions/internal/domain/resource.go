package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/kubewarden/policy-sdk-go/pkg/capabilities"
	"github.com/kubewarden/policy-sdk-go/pkg/capabilities/kubernetes"
)

// ResourceRequestValidator validates an incoming Resource Request.
type ResourceRequestValidator struct {
	migRegex *regexp.Regexp
}

func NewResourceRequestValidator() ResourceRequestValidator {
	re := regexp.MustCompile(`nvidia\.com/mig-.*`)
	return ResourceRequestValidator{
		migRegex: re,
	}
}

// IsAllowed verifies if a Pod's resource request is allowed.
// We look for the namespace's ResourceQuota to validate that a MIG Partition is allowed.
//
// Restrictions
//   - If there is no ResourceQuota and the resource is a MIG Partition, deny.
//   - If there is a ResourceQuota and the MIG Partition is not in it, deny.
//   - If there is a ResourceQuota and the MIG Partition is in it, allow.
func (v *ResourceRequestValidator) IsAllowed(
	_ context.Context,
	host *capabilities.Host,
	namespace, resource string,
) bool {
	// If the resource is not a MIG Partition, we can skip it.
	// We only want to validate MIG Partitions.
	if !v.isMigPartition(resource) {
		return true
	}

	// Try to get the namespace's ResourceQuota from Kubernetes.
	// If we cannot find it, deny the request.
	resourceQuotaList, err := v.findResourceQuotasByNamespace(host, namespace)
	if err != nil {
		return false
	}

	for _, resourceQuota := range resourceQuotaList.Items {
		for k := range resourceQuota.Spec.Hard {
			// in a ResourceQuota, the mig partition will have the prefix "requests."
			// for example, nvidia.com/mig-2g.24gb will be requests.nvidia.com/mig-2g.24gb
			if k == ("requests." + resource) {
				return true
			}
		}
	}

	// If we didn't find the MIG Partition in the ResourceQuota, then it should be denied.
	return false
}

// findResourceQuotasByNamespace asks Kubernetes for the namespace's ResourceQuota.
//
// We may have more than one ResourceQuota, so we need to be sure to check all of them.
func (v *ResourceRequestValidator) findResourceQuotasByNamespace(
	host *capabilities.Host,
	namespace string,
) (ResourceQuotaList, error) {
	kubeRequest := kubernetes.ListResourcesByNamespaceRequest{
		APIVersion: "v1",
		Kind:       "ResourceQuota",
		Namespace:  namespace,
	}

	response, err := kubernetes.ListResourcesByNamespace(host, kubeRequest)
	if err != nil {
		return ResourceQuotaList{}, err
	}

	list := ResourceQuotaList{}
	err = json.Unmarshal(response, &list)
	if err != nil {
		return ResourceQuotaList{}, fmt.Errorf("cannot unmarshall response into ResourceQutoaList: %w", err)
	}

	return list, nil
}

// isMigPartition checks with a regular expression whether a resource is a MIG Partition.
func (v *ResourceRequestValidator) isMigPartition(resource string) bool {
	return v.migRegex.Match([]byte(resource))
}
