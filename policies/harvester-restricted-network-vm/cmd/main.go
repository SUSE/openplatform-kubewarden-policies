package main

import (
	"context"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-restricted-network-vm/internal/logger"
	"github.com/wapc/wapc-guest-tinygo"
)

func main() {
	ctx := logger.ContextWithLogger(context.Background())

	wapc.RegisterFunctions(wapc.Functions{
		"validate": func(payload []byte) ([]byte, error) {
			return validate(ctx, payload)
		},
		"validate_settings": func(payload []byte) ([]byte, error) {
			return validateSettings(ctx, payload)
		},
	})
}
