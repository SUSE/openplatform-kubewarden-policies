package inbound

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-pci-devices/internal/core/logger"
	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-pci-devices/internal/domain"
	kubewarden "github.com/kubewarden/policy-sdk-go"
)

func ValidateSettings(ctx context.Context, payload []byte) ([]byte, error) {
	l := logger.FromContext(ctx)
	l.Info("validating settings")

	settings := domain.Settings{}
	err := json.Unmarshal(payload, &settings)
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message(fmt.Sprintf("Invalid settings JSON: %v", err)))
	}

	valid := settings.Valid(ctx)
	if !valid {
		return kubewarden.RejectSettings("settings are not valid")
	}

	return kubewarden.AcceptSettings()
}
