package inbound

import (
	"context"

	"github.com/kubewarden/policy-sdk-go/pkg/capabilities"
)

type resourceValidator interface {
	IsAllowed(ctx context.Context, host *capabilities.Host, namespace, resource string) bool
}
