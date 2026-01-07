package main

import (
	"context"
	"encoding/json"
	"testing"

	kubewardenProtocol "github.com/kubewarden/policy-sdk-go/protocol"
	kubewardenTesting "github.com/kubewarden/policy-sdk-go/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getVMObject(vmName, namespace, network string) virtualMachine {
	return virtualMachine{
		Metadata: vmMetadata{
			Namespace: namespace,
			Name:      vmName,
		},
		Spec: vmPayloadSpec{
			Template: vmTemplate{
				Spec: vmNetworks{
					Networks: []vmNetwork{
						{
							Multus: multus{
								NetworkName: network,
							},
						},
					},
				},
			},
		},
	}
}

func TestApproval(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		getPayload   func() []byte
		result       bool
		errorMessage string
		errorCode    uint16
	}{
		{
			name: "Approve: Empty settings",
			getPayload: func() []byte {
				settings := Settings{}
				vmObject := getVMObject("test-VM", "random-namespace", "random-network")
				payload, err := kubewardenTesting.BuildValidationRequest(&vmObject, &settings)
				assert.NoError(t, err)
				return payload
			},
			result: true,
		},
		{
			name: "Approve: allowed VM",
			getPayload: func() []byte {
				settings := Settings{
					NamespaceNetworkBindings: []NamespaceNetworkBinding{
						{
							Network:   "restricted-network",
							Namespace: "restricted-namespace",
						},
					},
				}

				vmObject := getVMObject("test-VM", "restricted-namespace", "restricted-network")

				payload, err := kubewardenTesting.BuildValidationRequest(&vmObject, &settings)
				assert.NoError(t, err)
				return payload
			},
			result: true,
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
			errorCode:    httpBadRequestStatusCode,
		},
		{
			name: "Reject: Bad settings",
			getPayload: func() []byte {
				return []byte(`
				{
					"settings": 123
				}`)
			},
			result:       false,
			errorMessage: "json: cannot unmarshal number into Go value of type main.Settings",
			errorCode:    httpBadRequestStatusCode,
		},
		{
			name: "Reject: Bad VM json",
			getPayload: func() []byte {
				settings := Settings{
					NamespaceNetworkBindings: []NamespaceNetworkBinding{
						{
							Network:   "restricted-network",
							Namespace: "restricted-namespace",
						},
					},
				}

				vmObject := "fake vm"

				payload, err := kubewardenTesting.BuildValidationRequest(&vmObject, &settings)
				assert.NoError(t, err)
				return payload
			},
			result:       false,
			errorMessage: "json: cannot unmarshal string into Go value of type main.virtualMachine",
			errorCode:    httpBadRequestStatusCode,
		},
		{
			name: "Reject: random network for a restricted namespace",
			getPayload: func() []byte {
				settings := Settings{
					NamespaceNetworkBindings: []NamespaceNetworkBinding{
						{
							Network:   "restricted-network",
							Namespace: "restricted-namespace",
						},
					},
				}

				vmObject := getVMObject("test-VM", "restricted-namespace", "random-network")
				payload, err := kubewardenTesting.BuildValidationRequest(&vmObject, &settings)
				assert.NoError(t, err)
				return payload
			},
			result:       false,
			errorMessage: "Network 'random-network' is not allowed for namespace: 'restricted-namespace'",
			errorCode:    httpBadRequestStatusCode,
		},
		{
			name: "Reject: multiple network interfaces, one from different network, one from restricted-network",
			getPayload: func() []byte {
				settings := Settings{
					NamespaceNetworkBindings: []NamespaceNetworkBinding{
						{
							Network:   "restricted-network",
							Namespace: "restricted-namespace",
						},
					},
				}
				vmObject := virtualMachine{
					Metadata: vmMetadata{
						Namespace: "restricted-namespace",
						Name:      "test-VM",
					},
					Spec: vmPayloadSpec{
						Template: vmTemplate{
							Spec: vmNetworks{
								Networks: []vmNetwork{
									{
										Multus: multus{
											NetworkName: "random-network",
										},
									},
									{
										Multus: multus{
											NetworkName: "restricted-network",
										},
									},
								},
							},
						},
					},
				}
				payload, err := kubewardenTesting.BuildValidationRequest(&vmObject, &settings)
				require.NoError(t, err)
				return payload
			},
			result:       false,
			errorMessage: "Network 'random-network' is not allowed for namespace: 'restricted-namespace'",
			errorCode:    httpBadRequestStatusCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := tt.getPayload()
			responsePayload, err := validate(ctx, payload)
			require.NoError(t, err)

			var response kubewardenProtocol.ValidationResponse
			err = json.Unmarshal(responsePayload, &response)
			require.NoError(t, err)

			assert.Equal(t, tt.result, response.Accepted)
			if !tt.result {
				assert.Equal(t, tt.errorMessage, *response.Message)
				assert.Equal(t, tt.errorCode, *response.Code)
			}
		})
	}
}
