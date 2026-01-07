package main

import (
	"context"
	"fmt"

	kubewarden "github.com/kubewarden/policy-sdk-go"
)

func validateSettings(ctx context.Context, payload []byte) ([]byte, error) {
	settings, err := parseSettings(payload)
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message(fmt.Sprintf("Invalid settings JSON: %v", err)))
	}

	valid := settings.IsValid(ctx)
	if !valid {
		return kubewarden.RejectSettings("Invalid settings")
	}

	return kubewarden.AcceptSettings()
}
