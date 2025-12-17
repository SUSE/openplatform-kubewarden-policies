package inbound

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/rke2-mig-partitions/internal/domain"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	"github.com/kubewarden/policy-sdk-go/pkg/capabilities"
	kubewardenProtocol "github.com/kubewarden/policy-sdk-go/protocol"
)

const HTTPBadRequestStatusCode = 400

func ValidateRequest(ctx context.Context, payload []byte, validator resourceValidator) ([]byte, error) {
	validationRequest := kubewardenProtocol.ValidationRequest{}
	err := json.Unmarshal(payload, &validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(HTTPBadRequestStatusCode))
	}

	host := capabilities.NewHost()

	podJSON := validationRequest.Request.Object

	podObject := domain.Pod{}
	err = json.Unmarshal(podJSON, &podObject)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(HTTPBadRequestStatusCode))
	}

	namespace := podObject.Metadata.Namespace

	for _, container := range podObject.Spec.Containers {
		for resourceRequest := range container.Resources.Requests {
			if !validator.IsAllowed(ctx, &host, namespace, resourceRequest) {
				return kubewarden.RejectRequest(
					kubewarden.Message(
						fmt.Sprintf(
							"MIG Partition '%s' is not allowed for namespace: '%s'", resourceRequest, namespace,
						)),
					kubewarden.Code(HTTPBadRequestStatusCode))
			}
		}
	}

	return kubewarden.AcceptRequest()
}
