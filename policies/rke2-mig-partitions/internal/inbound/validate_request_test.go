package inbound

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/rke2-mig-partitions/internal/domain"
	kubewardenProtocol "github.com/kubewarden/policy-sdk-go/protocol"
	kubewardenTesting "github.com/kubewarden/policy-sdk-go/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func getPodWithMigPartition(name, namespace, migPartition string) domain.Pod {
	return domain.Pod{
		Metadata: domain.Metadata{
			Namespace: namespace,
			Name:      name,
		},
		Spec: domain.PodSpec{
			Containers: []domain.ContainerSpec{
				{
					Resources: domain.PodSpecResources{
						Requests: map[string]interface{}{
							migPartition: 1,
						},
						Limits: map[string]interface{}{
							migPartition: 1,
						},
					},
				},
			},
		},
	}
}

func getPod(name, namespace string) domain.Pod {
	return domain.Pod{
		Metadata: domain.Metadata{
			Namespace: namespace,
			Name:      name,
		},
		Spec: domain.PodSpec{
			Containers: []domain.ContainerSpec{
				{
					Resources: domain.PodSpecResources{},
				},
			},
		},
	}
}

func TestApproval(t *testing.T) {
	ctx := context.Background()
	settings := domain.Settings{}

	tt := []struct {
		name          string
		getPayload    func() []byte
		shouldApprove bool
		result        bool
		errorMessage  string
		errorCode     uint16
	}{
		{
			name:          "Approve: No Mig Partition",
			shouldApprove: true,
			getPayload: func() []byte {
				pod := getPod("test", "random-namespace")
				payload, err := kubewardenTesting.BuildValidationRequest(&pod, &settings)
				assert.NoError(t, err)
				return payload
			},
			result: true,
		},
		{
			name:          "Approve with mig partition",
			shouldApprove: true,
			getPayload: func() []byte {
				vmObject := getPodWithMigPartition("test", "random-namespace", "nvidia.com/mig-1g.12gb")
				payload, err := kubewardenTesting.BuildValidationRequest(&vmObject, &settings)
				assert.NoError(t, err)
				return payload
			},
			result: true,
		},
		{
			name:          "Deny with mig partition",
			shouldApprove: false,
			getPayload: func() []byte {
				vmObject := getPodWithMigPartition("test", "random-namespace", "nvidia.com/mig-1g.12gb")
				payload, err := kubewardenTesting.BuildValidationRequest(&vmObject, &settings)
				assert.NoError(t, err)
				return payload
			},
			result:       false,
			errorMessage: "MIG Partition 'nvidia.com/mig-1g.12gb' is not allowed for namespace: 'random-namespace'",
			errorCode:    HTTPBadRequestStatusCode,
		},
		{
			name: "Reject: Bad payload",
			getPayload: func() []byte {
				return []byte(`
					{
						"request": test
					}`)
			},
			result:       false,
			errorMessage: "invalid character 'e' in literal true (expecting 'r')",
			errorCode:    HTTPBadRequestStatusCode,
		},
		{
			name: "Reject: Bad VM json",
			getPayload: func() []byte {
				vmObject := "fake vm"
				payload, err := kubewardenTesting.BuildValidationRequest(&vmObject, &settings)
				assert.NoError(t, err)
				return payload
			},
			result:       false,
			errorMessage: "json: cannot unmarshal string into Go value of type domain.Pod",
			errorCode:    HTTPBadRequestStatusCode,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			payload := tc.getPayload()
			validator := new(mockResourceValidator)
			validator.On("IsAllowed", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(tc.shouldApprove)
			responsePayload, err := ValidateRequest(ctx, payload, validator)
			require.NoError(t, err)

			var response kubewardenProtocol.ValidationResponse
			err = json.Unmarshal(responsePayload, &response)
			require.NoError(t, err)

			assert.Equal(t, tc.result, response.Accepted)
			if !tc.result {
				assert.Equal(t, tc.errorMessage, *response.Message)
				assert.Equal(t, tc.errorCode, *response.Code)
			}
		})
	}
}
