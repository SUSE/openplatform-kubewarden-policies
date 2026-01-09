package inbound

import (
	"context"

	kubewarden "github.com/kubewarden/policy-sdk-go"
)

func ValidateSettings(_ context.Context, _ []byte) ([]byte, error) {
	return kubewarden.AcceptSettings()
}
