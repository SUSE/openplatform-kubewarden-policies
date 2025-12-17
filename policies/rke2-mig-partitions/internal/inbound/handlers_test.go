package inbound

import (
	"context"

	"github.com/kubewarden/policy-sdk-go/pkg/capabilities"
	"github.com/stretchr/testify/mock"
)

var _ resourceValidator = &mockResourceValidator{}

type mockResourceValidator struct {
	mock.Mock
}

func (m *mockResourceValidator) IsAllowed(
	ctx context.Context,
	host *capabilities.Host,
	namespace, resource string,
) bool {
	args := m.Called(ctx, host, namespace, resource)

	return args.Get(0).(bool)
}
