package inbound_test

import (
	"context"
	"testing"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-pci-devices/internal/inbound"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

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
				"namespaceDeviceBindings": [
				{
					"namespace": "test-namespace-1",
					"device": "gpu-1"
				}
				]
			}`),
			result: `{"valid":true}`,
		},
		{
			name: "Invalid settings json",
			payload: []byte(`{
				"namespaceDeviceBindings": [
			}`),
			result: `{"valid":false,"message":"Invalid settings JSON: invalid character '}' looking for beginning of value"}`,
		},
		{
			name: "Missing namespace",
			payload: []byte(`{
				"namespaceDeviceBindings": [
					{
						"device": "gpu-1"
					}
				]
			}`),
			result: `{"valid":false,"message":"settings are not valid"}`,
		},
		{
			name: "Missing namespace",
			payload: []byte(`{
				"namespaceDeviceBindings": [
					{
						"namespace": "test-namespace-1"
					}
				]
			}`),
			result: `{"valid":false,"message":"settings are not valid"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := inbound.ValidateSettings(ctx, tt.payload)
			require.NoError(t, err)
			assert.Equal(t, tt.result, string(result))
		})
	}
}
