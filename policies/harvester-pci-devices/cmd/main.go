package main

import (
	"context"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-pci-devices/internal/core/logger"
	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-pci-devices/internal/inbound"
	"github.com/wapc/wapc-guest-tinygo"
)

func main() {
	ctx := logger.ContextWithLogger(context.Background())

	wapc.RegisterFunctions(wapc.Functions{
		"validate": func(payload []byte) ([]byte, error) {
			return inbound.ValidateRequest(ctx, payload)
		},
		"validate_settings": func(payload []byte) ([]byte, error) {
			return inbound.ValidateSettings(ctx, payload)
		},
	})
}
