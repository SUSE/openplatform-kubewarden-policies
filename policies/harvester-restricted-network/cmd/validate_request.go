package main

import (
	"context"
	"encoding/json"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-restricted-network/internal"
	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-restricted-network/internal/logger"
	"github.com/francoispqt/onelog"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewardenProtocol "github.com/kubewarden/policy-sdk-go/protocol"
)

func validateRequest(ctx context.Context, payload []byte) ([]byte, error) {
	validationRequest := kubewardenProtocol.ValidationRequest{}
	err := json.Unmarshal(payload, &validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(httpBadRequestStatusCode))
	}

	settings, err := parseSettings(validationRequest.Settings)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(httpBadRequestStatusCode))
	}

	networkAttachmentDefinition := internal.NetworkAttachmentDefinition{}
	err = json.Unmarshal(validationRequest.Request.Object, &networkAttachmentDefinition)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(httpBadRequestStatusCode))
	}

	namespace := networkAttachmentDefinition.Metadata.Namespace
	vlan := networkAttachmentDefinition.Spec.Config.VLAN
	l := logger.FromContext(ctx).With(func(entry onelog.Entry) {
		entry.String("namespace", namespace)
		entry.Int("vlan", vlan)
	})

	l.Info("NETWORK_CHECK namespace")
	if !settings.IsVLANAllowed(ctx, namespace, vlan) {
		l.Info("NETWORK_REJECTED namespace")
		return kubewarden.RejectRequest("Invalid request", kubewarden.Code(httpBadRequestStatusCode))
	}
	l.Info("NETWORK_ALLOWED namespace")

	return kubewarden.AcceptRequest()
}
