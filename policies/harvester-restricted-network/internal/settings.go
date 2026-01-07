package internal

import (
	"context"
	"fmt"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-restricted-network/internal/logger"
	"github.com/francoispqt/onelog"
)

type NamespaceVLANBinding struct {
	Namespace string `json:"namespace"`
	VLAN      int    `json:"vlan"`
}

// Settings is the structure that describes the policy settings.
type Settings struct {
	NamespaceVLANBindings []NamespaceVLANBinding `json:"namespaceVLANBindings"`
}

// IsValid verifies the Settings object, by ensuring that the namespace and VLANs are always defined.
func (s *Settings) IsValid(ctx context.Context) bool {
	l := logger.FromContext(ctx)

	for _, ns := range s.NamespaceVLANBindings {
		// Check if namespace and network are not empty
		if ns.Namespace == "" || ns.VLAN == 0 {
			l.DebugWithFields("Namespace or VLAN required", func(e onelog.Entry) {
				e.String("Namespace", ns.Namespace)
				e.Int("VLAN", ns.VLAN)
			})
			return false
		}
	}
	return true
}

// IsVLANAllowed verifies if a (namespace, VLAN) combination is allowed.
//
// Restrictions
//   - namespaces with vlan bindings can only accept a vlan bound to it
//   - VLANs with namespace bindings can only accept a namespace bound to it
//   - If a namespace and VLAN don't have a binding, then it's unrestricted
//
// example:
//
//	  settings:
//			- namespace: restricted-namespace
//			  network: 42
//
// Allowed:
//
//	{"namespace": "restricted-namespace", "vlan": "42"}
//	{"namespace": "random-namespace", "network": "1337"}
//
// Denied:
//
//	{"namespace": "restricted-namespace", "network": "1337"}
//	{"namespace": "random-namespace", "network": "42"}
func (s *Settings) IsVLANAllowed(ctx context.Context, namespace string, vlan int) bool {
	allowed := true
	l := logger.FromContext(ctx).With(func(e onelog.Entry) {
		e.String("Namespace", namespace)
		e.Int("VLAN", vlan)
	})

	for _, ns := range s.NamespaceVLANBindings {
		// if a namespace is bound, then its vlan must be bound to it
		if ns.Namespace == namespace && ns.VLAN != vlan {
			allowed = false
		}

		// if a vlan is bound, then its namespace must be bound to it
		if ns.VLAN == vlan && ns.Namespace != namespace {
			allowed = false
		}

		// vlan and namespace are bound
		if ns.VLAN == vlan && ns.Namespace == namespace {
			l.Debug("vlan and namespace matched")
			return true
		}
	}

	// if allowed is "true", it's because the namespace and vlan are not bound and are considered unrestricted
	if allowed {
		l.Debug(fmt.Sprintf("namespace '%s' and vlan '%d' are unrestricted", namespace, vlan))
		return true
	}

	l.Debug(fmt.Sprintf("namespace `%s` or vlan '%d' are restricted and not bound together", namespace, vlan))
	return false
}
