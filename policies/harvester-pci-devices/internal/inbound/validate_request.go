package inbound

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-pci-devices/internal/core/logger"
	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-pci-devices/internal/domain"
	"github.com/francoispqt/onelog"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewardenProtocol "github.com/kubewarden/policy-sdk-go/protocol"
)

const HTTPBadRequestStatusCode = 400

func ValidateRequest(ctx context.Context, payload []byte) ([]byte, error) {
	validationRequest := kubewardenProtocol.ValidationRequest{}
	err := json.Unmarshal(payload, &validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(HTTPBadRequestStatusCode))
	}

	settings, err := domain.NewSettingsFromValidationReq(&validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(HTTPBadRequestStatusCode))
	}

	virtualMachineJSON := validationRequest.Request.Object

	virtualMachineObject := domain.VirtualMachine{}
	err = json.Unmarshal(virtualMachineJSON, &virtualMachineObject)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(HTTPBadRequestStatusCode))
	}
	namespace := virtualMachineObject.Metadata.Namespace
	gpuList := virtualMachineObject.Spec.Template.Spec.Domain.Devices.GPUS
	pciDeviceList := virtualMachineObject.Spec.Template.Spec.Domain.Devices.HostDevices

	l := logger.FromContext(ctx).With(func(entry onelog.Entry) {
		entry.String("namespace", namespace)
		entry.String("devices", fmt.Sprintf("%+v", pciDeviceList))
	})

	l.Info("VM_CHECK namespace/device")
	for _, gpu := range append(gpuList, pciDeviceList...) {
		gpuName := gpu.Name
		if !settings.IsGPUAllowed(ctx, namespace, gpuName) {
			l.InfoWithFields("VM_REJECTED namespace/device", func(entry onelog.Entry) {
				entry.String("device", gpuName)
			})
			return kubewarden.RejectRequest(
				kubewarden.Message(
					fmt.Sprintf("PCI DEVICE '%s' is not allowed for namespace: '%s'", gpuName, namespace)),
				kubewarden.Code(HTTPBadRequestStatusCode))
		}
	}

	l.Info("VM_ALLOWED namespace")
	return kubewarden.AcceptRequest()
}
