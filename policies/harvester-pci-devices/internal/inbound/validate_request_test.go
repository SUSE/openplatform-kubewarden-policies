package inbound_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-pci-devices/internal/inbound"

	"github.com/stretchr/testify/require"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-pci-devices/internal/domain"
	kubewardenProtocol "github.com/kubewarden/policy-sdk-go/protocol"
	kubewardenTesting "github.com/kubewarden/policy-sdk-go/testing"
	"github.com/stretchr/testify/assert"
)

func getVMObjectGPU(vmName, namespace, gpu string) domain.VirtualMachine {
	return domain.VirtualMachine{
		Metadata: domain.Metadata{
			Namespace: namespace,
			Name:      vmName,
		},
		Spec: domain.VirtualMachineSpec{
			Template: domain.VirtualMachineSpecTemplate{
				Spec: domain.VirtualMachineTemplateSpec{
					Domain: domain.Domain{
						Devices: domain.Devices{
							GPUS: []domain.PCIDevice{
								{
									Name:       gpu,
									DeviceName: "gpu-device-name",
								},
							},
						},
					},
				},
			},
		},
	}
}

func getVMObjectPCI(vmName, namespace, gpu string) domain.VirtualMachine {
	return domain.VirtualMachine{
		Metadata: domain.Metadata{
			Namespace: namespace,
			Name:      vmName,
		},
		Spec: domain.VirtualMachineSpec{
			Template: domain.VirtualMachineSpecTemplate{
				Spec: domain.VirtualMachineTemplateSpec{
					Domain: domain.Domain{
						Devices: domain.Devices{
							HostDevices: []domain.PCIDevice{
								{
									Name:       gpu,
									DeviceName: "gpu-device-name",
								},
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
			name: "Deny: Empty settings (GPU)",
			getPayload: func() []byte {
				settings := domain.Settings{}
				vmObject := getVMObjectGPU("test-VM", "random-namespace", "random-gpu")
				payload, err := kubewardenTesting.BuildValidationRequest(&vmObject, &settings)
				assert.NoError(t, err)
				return payload
			},
			errorMessage: "PCI DEVICE 'random-gpu' is not allowed for namespace: 'random-namespace'",
			errorCode:    inbound.HTTPBadRequestStatusCode,
			result:       false,
		},
		{
			name: "Deny: Empty settings (PCI)",
			getPayload: func() []byte {
				settings := domain.Settings{}
				vmObject := getVMObjectPCI("test-VM", "random-namespace", "random-gpu")
				payload, err := kubewardenTesting.BuildValidationRequest(&vmObject, &settings)
				assert.NoError(t, err)
				return payload
			},
			errorMessage: "PCI DEVICE 'random-gpu' is not allowed for namespace: 'random-namespace'",
			errorCode:    inbound.HTTPBadRequestStatusCode,
			result:       false,
		},
		{
			name: "Approve: bound namespace and Device",
			getPayload: func() []byte {
				settings := domain.Settings{
					NamespaceDeviceBindings: []domain.NamespaceDeviceBinding{
						{
							Device:    "gpu-1",
							Namespace: "namespace-1",
						},
					},
				}

				vmObject := getVMObjectGPU("test-VM", "namespace-1", "gpu-1")

				payload, err := kubewardenTesting.BuildValidationRequest(&vmObject, &settings)
				assert.NoError(t, err)
				return payload
			},
			result: true,
		},
		{
			name: "Approve: bound namespace and PCI Device",
			getPayload: func() []byte {
				settings := domain.Settings{
					NamespaceDeviceBindings: []domain.NamespaceDeviceBinding{
						{
							Device:    "gpu-1",
							Namespace: "namespace-1",
						},
					},
				}

				vmObject := getVMObjectPCI("test-VM", "namespace-1", "gpu-1")

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
			errorCode:    inbound.HTTPBadRequestStatusCode,
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
			errorMessage: "json: cannot unmarshal number into Go value of type domain.Settings",
			errorCode:    inbound.HTTPBadRequestStatusCode,
		},
		{
			name: "Reject: Bad VM json",
			getPayload: func() []byte {
				settings := domain.Settings{
					NamespaceDeviceBindings: []domain.NamespaceDeviceBinding{
						{
							Device:    "gpu-1",
							Namespace: "namespace-2",
						},
					},
				}

				vmObject := "fake vm"

				payload, err := kubewardenTesting.BuildValidationRequest(&vmObject, &settings)
				assert.NoError(t, err)
				return payload
			},
			result:       false,
			errorMessage: "json: cannot unmarshal string into Go value of type domain.VirtualMachine",
			errorCode:    inbound.HTTPBadRequestStatusCode,
		},
		{
			name: "Reject: random Device for a bound namespace",
			getPayload: func() []byte {
				settings := domain.Settings{
					NamespaceDeviceBindings: []domain.NamespaceDeviceBinding{
						{
							Device:    "gpu-random",
							Namespace: "namespace-1",
						},
					},
				}

				vmObject := getVMObjectGPU("test-VM", "namespace-1", "gpu-1")
				payload, err := kubewardenTesting.BuildValidationRequest(&vmObject, &settings)
				assert.NoError(t, err)
				return payload
			},
			result:       false,
			errorMessage: "PCI DEVICE 'gpu-1' is not allowed for namespace: 'namespace-1'",
			errorCode:    inbound.HTTPBadRequestStatusCode,
		},
		{
			name: "Reject: random PCI DEVICE for a bound namespace",
			getPayload: func() []byte {
				settings := domain.Settings{
					NamespaceDeviceBindings: []domain.NamespaceDeviceBinding{
						{
							Device:    "gpu-random",
							Namespace: "namespace-1",
						},
					},
				}

				vmObject := getVMObjectPCI("test-VM", "namespace-1", "gpu-1")
				payload, err := kubewardenTesting.BuildValidationRequest(&vmObject, &settings)
				assert.NoError(t, err)
				return payload
			},
			result:       false,
			errorMessage: "PCI DEVICE 'gpu-1' is not allowed for namespace: 'namespace-1'",
			errorCode:    inbound.HTTPBadRequestStatusCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := tt.getPayload()
			responsePayload, err := inbound.ValidateRequest(ctx, payload)
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
