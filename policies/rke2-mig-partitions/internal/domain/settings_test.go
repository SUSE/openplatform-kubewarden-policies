package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	kubewardenProtocol "github.com/kubewarden/policy-sdk-go/protocol"
	"github.com/stretchr/testify/assert"
)

func TestSettings_Valid(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		name         string
		settings     Settings
		expectResult bool
	}{
		{
			name:         "Valid settings",
			settings:     Settings{},
			expectResult: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.settings.Valid(ctx)

			assert.Equal(t, tc.expectResult, result)
		})
	}
}

func TestNewSettingsFromValidationReq(t *testing.T) {
	settingsJSON := []byte(`{}`)
	validationRequest := &kubewardenProtocol.ValidationRequest{
		Request:  kubewardenProtocol.KubernetesAdmissionRequest{},
		Settings: settingsJSON,
	}

	_, err := NewSettingsFromValidationReq(validationRequest)
	require.NoError(t, err)
}
