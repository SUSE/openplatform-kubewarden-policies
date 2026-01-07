package domain

import (
	"context"
	"encoding/json"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-pci-devices/internal/core/logger"
	"github.com/francoispqt/onelog"
	kubewardenProtocol "github.com/kubewarden/policy-sdk-go/protocol"
)

type NamespaceDeviceBinding struct {
	Namespace string `json:"namespace"`
	Device    string `json:"device"`
}

// Settings is the structure that describes the policy settings.
type Settings struct {
	NamespaceDeviceBindings []NamespaceDeviceBinding `json:"namespaceDeviceBindings"`
}

func NewSettingsFromValidationReq(validationReq *kubewardenProtocol.ValidationRequest) (Settings, error) {
	settings := Settings{}
	err := json.Unmarshal(validationReq.Settings, &settings)
	return settings, err
}

func (s *Settings) Valid(ctx context.Context) bool {
	l := logger.FromContext(ctx)

	for _, binding := range s.NamespaceDeviceBindings {
		// Check if namespace and device are not empty
		if binding.Namespace == "" || binding.Device == "" {
			l.InfoWithFields("invalid settings, namespace and device must be specified", func(entry onelog.Entry) {
				entry.String("namespace", binding.Namespace)
				entry.String("device", binding.Device)
			})

			return false
		}
	}

	return true
}

// IsGPUAllowed verifies if a (namespace, device) combination is allowed.
//
// Restrictions
//   - namespaces with pci device bindings can only accept a pci device bound to it
//   - pci devices with namespace bindings can only accept a namespace bound to it
//   - If a namespace and a pci device don't have a binding, then it's restricted
//
// example:
//
//	  settings:
//			- namespace: namespace-01
//			  device: tekton27a-000001010
//
// Allowed:
//
//	{"namespace": "namespace-01", "device": "tekton27a-000001010"}
//
// Denied:
//
//	{"namespace": "namespace-02", "device": "tekton27a-000001010"}
//	{"namespace": "random-namespace", "device": "tekton28a-000001010"}
func (s *Settings) IsGPUAllowed(ctx context.Context, namespace, device string) bool {
	l := logger.FromContext(ctx).With(func(entry onelog.Entry) {
		entry.String("namespace", namespace)
		entry.String("device", device)
	})

	for _, ns := range s.NamespaceDeviceBindings {
		// device and namespace are bound
		if ns.Device == device && ns.Namespace == namespace {
			l.Debug("device and namespace matched")
			return true
		}
	}

	// if allowed is "true", it's because the namespace and device are not bound and are considered unrestricted
	l.Debug("Device is restricted")
	return false
}
