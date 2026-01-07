package domain_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-pci-devices/internal/domain"

	"github.com/stretchr/testify/require"

	kubewardenProtocol "github.com/kubewarden/policy-sdk-go/protocol"
	"github.com/stretchr/testify/assert"
)

func TestSettings_Valid(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		name         string
		settings     domain.Settings
		expectResult bool
	}{
		{
			name: "Valid settings with all namespace and Device values",
			settings: domain.Settings{
				NamespaceDeviceBindings: []domain.NamespaceDeviceBinding{
					{Namespace: "namespace1", Device: "device-1"},
					{Namespace: "namespace2", Device: "device-2"},
				},
			},
			expectResult: true,
		},
		{
			name: "Invalid settings with empty namespace",
			settings: domain.Settings{
				NamespaceDeviceBindings: []domain.NamespaceDeviceBinding{
					{Namespace: "", Device: "device-1"},
				},
			},
			expectResult: false,
		},
		{
			name: "Invalid settings with empty Device",
			settings: domain.Settings{
				NamespaceDeviceBindings: []domain.NamespaceDeviceBinding{
					{Namespace: "namespace1", Device: ""},
				},
			},
			expectResult: false,
		},
		{
			name: "Invalid settings with empty namespace and Device",
			settings: domain.Settings{
				NamespaceDeviceBindings: []domain.NamespaceDeviceBinding{
					{Namespace: "", Device: ""},
				},
			},
			expectResult: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.settings.Valid(ctx)

			assert.Equal(t, tc.expectResult, result)
		})
	}
}

func TestSettings_IsGPUAllowed(t *testing.T) {
	ctx := context.Background()

	settingsJSON := []byte(`{
		"namespaceDeviceBindings": [
		{
			"namespace": "test-namespace-1",
			"device": "test-device-1"
		},
		{
			"namespace": "test-namespace-2",
			"device": "test-device-2"
		},
		{
			"namespace": "test-namespace-3",
			"device": "test-device-2"
		}
		]
	}`)
	settings := domain.Settings{}
	err := json.Unmarshal(settingsJSON, &settings)
	require.NoError(t, err)

	tests := []struct {
		name      string
		namespace string
		device    string
		result    bool
	}{
		{
			name:      "Valid pair namespace and device",
			namespace: "test-namespace-1",
			device:    "test-device-1",
			result:    true,
		},
		{
			name:      "Valid pair: when the device is also allowed in a different namespace",
			namespace: "test-namespace-2",
			device:    "test-device-2",
			result:    true,
		},
		{
			name:      "Valid pair: when the device is from a different namespace",
			namespace: "test-namespace-3",
			device:    "test-device-2",
			result:    true,
		},
		{
			name:      "Invalid pair: namespace with an unbound Device",
			namespace: "random-namespace",
			device:    "random-device",
			result:    false,
		},
		{
			name:      "Invalid pair: namespace with a bound device",
			namespace: "test-namespace-1",
			device:    "invalid-device",
			result:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := settings.IsGPUAllowed(ctx, tt.namespace, tt.device)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestNewSettingsFromValidationReq(t *testing.T) {
	settingsJSON := []byte(`{
		"namespaceDeviceBindings": [
		{
			"namespace": "test-cc-namespace-1",
			"device": "test-cc-device-1"
		}
		]
	}`)
	validationRequest := &kubewardenProtocol.ValidationRequest{
		Request:  kubewardenProtocol.KubernetesAdmissionRequest{},
		Settings: settingsJSON,
	}

	newSettings, err := domain.NewSettingsFromValidationReq(validationRequest)
	require.NoError(t, err)
	assert.Equal(t, "test-cc-namespace-1", newSettings.NamespaceDeviceBindings[0].Namespace)
}
