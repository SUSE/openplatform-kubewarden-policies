package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValid(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		settings    Settings
		result      bool
		expectError bool
	}{
		{
			name: "Valid settings with all namespace and network values",
			settings: Settings{
				NamespaceVLANBindings: []NamespaceVLANBinding{
					{Namespace: "namespace1", VLAN: 42},
					{Namespace: "namespace2", VLAN: 1337},
				},
			},
			result:      true,
			expectError: false,
		},
		{
			name: "Invalid settings with empty namespace",
			settings: Settings{
				NamespaceVLANBindings: []NamespaceVLANBinding{
					{Namespace: "", VLAN: 42},
				},
			},
			result:      false,
			expectError: true,
		},
		{
			name: "Invalid settings with empty network",
			settings: Settings{
				NamespaceVLANBindings: []NamespaceVLANBinding{
					{Namespace: "namespace1", VLAN: 0},
				},
			},
			result:      false,
			expectError: true,
		},
		{
			name: "Invalid settings with empty namespace and network",
			settings: Settings{
				NamespaceVLANBindings: []NamespaceVLANBinding{
					{Namespace: "", VLAN: 0},
				},
			},
			result:      false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.settings.IsValid(ctx)

			assert.Equal(t, tt.result, result)
		})
	}
}

func TestIsNetworkAllowed(t *testing.T) {
	settings := Settings{
		NamespaceVLANBindings: []NamespaceVLANBinding{
			{Namespace: "test-restricted-namespace-1", VLAN: 42},
			{Namespace: "test-restricted-namespace-2", VLAN: 1337},
			{Namespace: "test-restricted-namespace-3", VLAN: 1338},
		},
	}
	ctx := context.Background()

	tests := []struct {
		name      string
		namespace string
		vlan      int
		result    bool
	}{
		{
			name:      "Valid pair test 1",
			namespace: "test-restricted-namespace-1",
			vlan:      42,
			result:    true,
		},
		{
			name:      "Valid pair test 2",
			namespace: "test-restricted-namespace-2",
			vlan:      1337,
			result:    true,
		},
		{
			name:      "Valid pair test 3",
			namespace: "test-restricted-namespace-3",
			vlan:      1338,
			result:    true,
		},
		{
			name:      "Valid pair: non-restricted-namespace with a non-restricted network",
			namespace: "random-non-restricted-namespace",
			vlan:      100,
			result:    true,
		},
		{
			name:      "Invalid pair: restricted-namespace with a non-restricted network",
			namespace: "test-restricted-namespace-1",
			vlan:      443,
			result:    false,
		},
		{
			name:      "Invalid pair: non-restricted-namespace with a restricted network",
			namespace: "random-non-restricted-namespace",
			vlan:      42,
			result:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := settings.IsVLANAllowed(ctx, tt.namespace, tt.vlan)
			assert.Equal(t, tt.result, result)
		})
	}
}
