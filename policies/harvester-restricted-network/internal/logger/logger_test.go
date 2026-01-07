package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	ctx := ContextWithLogger(context.Background())
	logger := FromContext(ctx)
	logger2 := FromContext(ctx) // they should be the same object

	assert.NotNil(t, logger)
	assert.NotNil(t, logger2)
	assert.Equal(t, logger, logger2)
}

func TestFromContext_NoLogger(t *testing.T) {
	ctx := context.Background()

	logger := FromContext(ctx)
	assert.NotNil(t, logger)
}
