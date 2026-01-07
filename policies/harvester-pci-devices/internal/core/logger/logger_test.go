package logger_test

import (
	"context"
	"testing"

	"github.com/SUSE/openplatform-kubewarden-policies/policies/harvester-pci-devices/internal/core/logger"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	ctx := logger.ContextWithLogger(context.Background())
	log := logger.FromContext(ctx)
	log2 := logger.FromContext(ctx) // they should be the same object

	assert.NotNil(t, log)
	assert.NotNil(t, log2)
	assert.Equal(t, log, log2)
}

func TestFromContext_NoLogger(t *testing.T) {
	ctx := context.Background()

	log := logger.FromContext(ctx)
	assert.NotNil(t, log)
}
