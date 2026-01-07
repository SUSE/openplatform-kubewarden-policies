package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-restricted-network-vm/internal/logger"
	"github.com/francoispqt/onelog"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewardenProtocol "github.com/kubewarden/policy-sdk-go/protocol"
)

type SettingsError struct {
	Namespace string
	Network   string
}

func (s *SettingsError) Error() string {
	return fmt.Sprintf(
		"Invalid settings error, namespace and network must be specified. Namespace: '%s', Network: '%s'",
		s.Namespace,
		s.Network,
	)
}

type NamespaceNetworkBinding struct {
	Namespace string `json:"namespace"`
	Network   string `json:"network"`
}

// Settings is the structure that describes the policy settings.
type Settings struct {
	NamespaceNetworkBindings []NamespaceNetworkBinding `json:"namespaceNetworkBindings"`
}

func NewSettingsFromValidationReq(validationReq *kubewardenProtocol.ValidationRequest) (Settings, error) {
	settings := Settings{}
	err := json.Unmarshal(validationReq.Settings, &settings)
	return settings, err
}

func (s *Settings) valid() (bool, error) {
	for _, ns := range s.NamespaceNetworkBindings {
		// Check if namespace and network are not empty
		if ns.Namespace == "" || ns.Network == "" {
			return false, &SettingsError{
				Namespace: ns.Namespace,
				Network:   ns.Network,
			}
		}
	}
	return true, nil
}

func validateSettings(ctx context.Context, payload []byte) ([]byte, error) {
	l := logger.FromContext(ctx)
	l.Info("validating settings")

	settings := Settings{}
	err := json.Unmarshal(payload, &settings)
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message(fmt.Sprintf("Invalid settings JSON: %v", err)))
	}

	valid, err := settings.valid()
	if !valid || err != nil {
		return kubewarden.RejectSettings("settings are not valid")
	}

	return kubewarden.AcceptSettings()
}

// isNetworkAllowed verifies if a (namespace, network) combination is allowed.
//
// Restrictions
//   - namespaces with network bindings can only accept a network bound to it
//   - networks with namespace bindings can only accept a namespace bound to it
//   - If a namespace and network don't have a binding, then it's unrestricted
//
// example:
//
//	  settings:
//			- namespace: restricted-namespace
//			  network: restricted-network
//
// Allowed:
//
//	{"namespace": "restricted-namespace", "network": "restricted-network"}
//	{"namespace": "random-namespace", "network": "random-network"}
//
// Denied:
//
//	{"namespace": "restricted-namespace", "network": "random-network"}
//	{"namespace": "random-namespace", "network": "restricted-network"}
func (s *Settings) isNetworkAllowed(ctx context.Context, namespace, network string) bool {
	l := logger.FromContext(ctx).With(func(entry onelog.Entry) {
		entry.String("namespace", namespace)
		entry.String("network", network)
	})
	allowed := true

	for _, ns := range s.NamespaceNetworkBindings {
		// if a namespace is bound, then its network must be bound to it
		if ns.Namespace == namespace && ns.Network != network {
			allowed = false
		}

		// if a network is bound, then its namespace must be bound to it
		if ns.Network == network && ns.Namespace != namespace {
			allowed = false
		}

		// network and namespace are bound
		if ns.Network == network && ns.Namespace == namespace {
			l.Debug("network and namespace matched")
			return true
		}
	}

	// if allowed is "true", it's because the namespace and network are not bound and are considered unrestricted
	if allowed {
		l.Debug("namespace and network are unrestricted")
		return true
	}

	l.Debug("namespace or network is restricted and cannot bound together")
	return false
}
