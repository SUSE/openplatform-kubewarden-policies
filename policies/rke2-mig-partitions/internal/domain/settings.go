package domain

import (
	"context"

	kubewardenProtocol "github.com/kubewarden/policy-sdk-go/protocol"
)

// Settings is the structure that describes the policy settings.
type Settings struct{}

func NewSettingsFromValidationReq(_ *kubewardenProtocol.ValidationRequest) (Settings, error) {
	return Settings{}, nil
}

func (s *Settings) Valid(_ context.Context) bool {
	return true
}
