package domain

import (
	"context"
	"testing"

	"github.com/kubewarden/policy-sdk-go/pkg/capabilities"
	"github.com/kubewarden/policy-sdk-go/pkg/capabilities/mocks"
	"github.com/stretchr/testify/assert"
)

func TestResourceRequestValidator_IsAllowed(t *testing.T) {
	ctx := context.Background()
	validator := NewResourceRequestValidator()
	expectedInputPayload := `{"api_version":"v1","kind":"ResourceQuota","namespace":"default"}`
	response12GBMig := `{"apiVersion":"v1","items":[{"apiVersion":"v1","kind":"ResourceQuota","metadata":{"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"v1\",\"kind\":\"ResourceQuota\",\"metadata\":{\"annotations\":{},\"name\":\"gpu-test\",\"namespace\":\"default\"},\"spec\":{\"hard\":{\"requests.nvidia.com/mig-1g.12gb\":\"1\"}}}\n"},"creationTimestamp":"2025-08-20T15:54:04Z","name":"gpu-test","namespace":"gpu-test","resourceVersion":"12135155","uid":"8ea9f464-5414-4d0e-b80a-8f83eb6dece9"},"spec":{"hard":{"requests.nvidia.com/mig-1g.12gb":"1"}},"status":{"hard":{"requests.nvidia.com/mig-1g.12gb":"1"},"used":{"requests.nvidia.com/mig-1g.12gb":"0"}}}],"kind":"List","metadata":{"resourceVersion":""}}`
	response24GBMig := `{"apiVersion":"v1","items":[{"apiVersion":"v1","kind":"ResourceQuota","metadata":{"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"v1\",\"kind\":\"ResourceQuota\",\"metadata\":{\"annotations\":{},\"name\":\"gpu-test\",\"namespace\":\"default\"},\"spec\":{\"hard\":{\"requests.nvidia.com/mig-2g.24gb\":\"1\"}}}\n"},"creationTimestamp":"2025-08-20T15:54:04Z","name":"gpu-test","namespace":"gpu-test","resourceVersion":"12135155","uid":"8ea9f464-5414-4d0e-b80a-8f83eb6dece9"},"spec":{"hard":{"requests.nvidia.com/mig-2g.24gb":"1"}},"status":{"hard":{"requests.nvidia.com/mig-2g.24gb":"1"},"used":{"requests.nvidia.com/mig-2g.24gb":"0"}}}],"kind":"List","metadata":{"resourceVersion":""}}`
	responseNoMig := `{"apiVersion":"v1","items":[],"kind":"List","metadata":{"resourceVersion":""}}`

	tt := []struct {
		name          string
		response      string
		responseError error
		resource      string
		result        bool
	}{
		{
			name:     "Valid 12gb mig request with a 12gb mig in its ResourceQuota",
			response: response12GBMig,
			resource: "nvidia.com/mig-1g.12gb",
			result:   true,
		},
		{
			name:     "Valid 24gb mig request with a 24gb mig in its ResourceQuota",
			response: response24GBMig,
			resource: "nvidia.com/mig-2g.24gb",
			result:   true,
		},
		{
			name:     "Valid not a mig partition",
			response: responseNoMig,
			resource: "foobar.com/barfoo",
			result:   true,
		},
		{
			name:     "Invalid 12gb mig request with a 24gb mig in its ResourceQuota",
			response: response24GBMig,
			resource: "nvidia.com/mig-1g.12gb",
			result:   false,
		},
		{
			name:     "Invalid 24gb mig request with a 12gb mig in its ResourceQuota",
			response: response12GBMig,
			resource: "nvidia.com/mig-2g.24gb",
			result:   false,
		},
		{
			name:          "Invalid ResourceQuota request failed",
			response:      response12GBMig,
			resource:      "nvidia.com/mig-1g.12gb",
			responseError: assert.AnError,
			result:        false,
		},
		{
			name:     "Invalid ResourceQuota request bad json",
			response: "foobar",
			resource: "nvidia.com/mig-1g.12gb",
			result:   false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockWapcClient := &mocks.MockWapcClient{}

			mockWapcClient.
				EXPECT().
				HostCall("kubewarden", "kubernetes", "list_resources_by_namespace", []byte(expectedInputPayload)).
				Return([]byte(tc.response), tc.responseError).
				Times(1)

			host := &capabilities.Host{
				Client: mockWapcClient,
			}

			result := validator.IsAllowed(ctx, host, "default", tc.resource)
			assert.Equal(t, tc.result, result)
		})
	}
}
