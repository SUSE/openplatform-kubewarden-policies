package main

import (
	"context"
	"encoding/json"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-restricted-network/internal"
	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-restricted-network/internal/logger"
	"github.com/wapc/wapc-guest-tinygo"
)

const httpBadRequestStatusCode = 400

func main() {
	ctx := logger.ContextWithLogger(context.Background())

	wapc.RegisterFunctions(wapc.Functions{
		"validate": func(payload []byte) ([]byte, error) {
			return validateRequest(ctx, payload)
		},
		"validate_settings": func(payload []byte) ([]byte, error) {
			return validateSettings(ctx, payload)
		},
	})
}

func parseSettings(payload []byte) (internal.Settings, error) {
	settings := internal.Settings{}
	err := json.Unmarshal(payload, &settings)
	if err != nil {
		return internal.Settings{}, err
	}

	return settings, nil
}
