package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-restricted-network-vm/internal/logger"

	"github.com/francoispqt/onelog"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewardenProtocol "github.com/kubewarden/policy-sdk-go/protocol"
)

const httpBadRequestStatusCode = 400

func validate(ctx context.Context, payload []byte) ([]byte, error) {
	validationRequest := kubewardenProtocol.ValidationRequest{}
	err := json.Unmarshal(payload, &validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(httpBadRequestStatusCode))
	}

	settings, err := NewSettingsFromValidationReq(&validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(httpBadRequestStatusCode))
	}

	virtualMachineJSON := validationRequest.Request.Object

	virtualMachineObject := virtualMachine{}
	err = json.Unmarshal(virtualMachineJSON, &virtualMachineObject)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(httpBadRequestStatusCode))
	}
	namespace := virtualMachineObject.Metadata.Namespace
	vmNetworkList := virtualMachineObject.Spec.Template.Spec.Networks

	l := logger.FromContext(ctx).With(func(entry onelog.Entry) {
		entry.String("namespace", namespace)
		entry.String("networks", fmt.Sprintf("%+v", vmNetworkList))
	})

	l.Info("VM_CHECK namespace/network")
	for _, network := range vmNetworkList {
		networkName := network.Multus.NetworkName
		if !settings.isNetworkAllowed(ctx, namespace, networkName) {
			l.InfoWithFields("VM_REJECTED namespace/network", func(entry onelog.Entry) {
				entry.String("network", networkName)
			})
			return kubewarden.RejectRequest(
				kubewarden.Message(
					fmt.Sprintf("Network '%s' is not allowed for namespace: '%s'", networkName, namespace)),
				kubewarden.Code(httpBadRequestStatusCode))
		}
	}

	l.Info("VM_ALLOWED namespace")
	return kubewarden.AcceptRequest()
}
