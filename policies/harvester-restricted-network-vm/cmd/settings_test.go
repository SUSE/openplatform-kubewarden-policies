package main

import (
	"context"
	"encoding/json"

	"testing"

	kubewardenProtocol "github.com/kubewarden/policy-sdk-go/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValid(t *testing.T) {
	tests := []struct {
		name          string
		settings      Settings
		result        bool
		expectError   bool
		expectedError error
	}{
		{
			name: "Valid settings with all namespace and network values",
			settings: Settings{
				NamespaceNetworkBindings: []NamespaceNetworkBinding{
					{Namespace: "namespace1", Network: "network1"},
					{Namespace: "namespace2", Network: "network2"},
				},
			},
			result:        true,
			expectError:   false,
			expectedError: nil,
		},
		{
			name: "Invalid settings with empty namespace",
			settings: Settings{
				NamespaceNetworkBindings: []NamespaceNetworkBinding{
					{Namespace: "", Network: "network1"},
				},
			},
			result:      false,
			expectError: true,
			expectedError: &SettingsError{
				Namespace: "",
				Network:   "network1",
			},
		},
		{
			name: "Invalid settings with empty network",
			settings: Settings{
				NamespaceNetworkBindings: []NamespaceNetworkBinding{
					{Namespace: "namespace1", Network: ""},
				},
			},
			result:      false,
			expectError: true,
			expectedError: &SettingsError{
				Namespace: "namespace1",
				Network:   "",
			},
		},
		{
			name: "Invalid settings with empty namespace and network",
			settings: Settings{
				NamespaceNetworkBindings: []NamespaceNetworkBinding{
					{Namespace: "", Network: ""},
				},
			},
			result:      false,
			expectError: true,
			expectedError: &SettingsError{
				Namespace: "",
				Network:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.settings.valid()

			if tt.expectError {
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.False(t, result)
			} else {
				require.NoError(t, err)
				assert.True(t, result)
			}
		})
	}
}

func TestValidateSettings(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		payload []byte
		result  string
	}{
		{
			name: "Valid settings",
			payload: []byte(`{
				"namespaceNetworkBindings": [
				{
					"namespace": "test-restricted-namespace-1",
					"network": "test-restricted-network-1"
				}
				]
			}`),
			result: `{"valid":true}`,
		},
		{
			name: "Invalid settings json",
			payload: []byte(`{
				"namespaceNetworkBindings": [
			}`),
			result: `{"valid":false,"message":"Invalid settings JSON: invalid character '}' looking for beginning of value"}`,
		},
		{
			name: "Missing namespace",
			payload: []byte(`{
				"namespaceNetworkBindings": [
					{
						"network": "test-restricted-network-1"
					}
				]
			}`),
			result: `{"valid":false,"message":"settings are not valid"}`,
		},
		{
			name: "Missing namespace",
			payload: []byte(`{
				"namespaceNetworkBindings": [
					{
						"namespace": "test-restricted-namespace-1"
					}
				]
			}`),
			result: `{"valid":false,"message":"settings are not valid"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateSettings(ctx, tt.payload)
			require.NoError(t, err)
			assert.Equal(t, tt.result, string(result))
		})
	}
}

func TestIsNetworkAllowed(t *testing.T) {
	ctx := context.Background()

	settingsJSON := []byte(`{
		"namespaceNetworkBindings": [
		{
			"namespace": "test-restricted-namespace-1",
			"network": "test-restricted-network-1"
		},
		{
			"namespace": "test-restricted-namespace-2",
			"network": "test-restricted-network-2"
		},
		{
			"namespace": "test-restricted-namespace-3",
			"network": "test-restricted-network-2"
		}
		]
	}`)
	settings := Settings{}
	err := json.Unmarshal(settingsJSON, &settings)
	require.NoError(t, err)

	tests := []struct {
		name      string
		namespace string
		network   string
		result    bool
	}{
		{
			name:      "Valid pair namespace and network",
			namespace: "test-restricted-namespace-1",
			network:   "test-restricted-network-1",
			result:    true,
		},
		{
			name:      "Valid pair: when the network is also allowed in a different namespace",
			namespace: "test-restricted-namespace-2",
			network:   "test-restricted-network-2",
			result:    true,
		},
		{
			name:      "Valid pair: when the network is from a different namespace",
			namespace: "test-restricted-namespace-3",
			network:   "test-restricted-network-2",
			result:    true,
		},
		{
			name:      "Valid pair: non-restricted-namespace with a non-restricted network",
			namespace: "random-non-restricted-namespace",
			network:   "random-non-restricted-network",
			result:    true,
		},
		{
			name:      "Invalid pair: restricted-namespace with a non-restricted network",
			namespace: "test-restricted-namespace-1",
			network:   "random-non-restricted-network",
			result:    false,
		},
		{
			name:      "Invalid pair: non-restricted-namespace with a restricted network",
			namespace: "random-non-restricted-namespace",
			network:   "test-restricted-network-2",
			result:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := settings.isNetworkAllowed(ctx, tt.namespace, tt.network)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestNewSettingsFromValidationReq(t *testing.T) {
	settingsJSON := []byte(`{
		"namespaceNetworkBindings": [
		{
			"namespace": "test-restricted-namespace-1",
			"network": "test-restricted-network-1"
		}
		]
	}`)
	validationRequest := &kubewardenProtocol.ValidationRequest{
		Request:  kubewardenProtocol.KubernetesAdmissionRequest{},
		Settings: settingsJSON,
	}

	newSettings, err := NewSettingsFromValidationReq(validationRequest)
	require.NoError(t, err)
	assert.Equal(t, "test-restricted-namespace-1", newSettings.NamespaceNetworkBindings[0].Namespace)
}
