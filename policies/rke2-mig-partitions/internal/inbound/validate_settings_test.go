package inbound

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestValidateSettings(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		name    string
		payload []byte
		result  string
	}{
		{
			name:    "Valid settings",
			payload: []byte(`{}`),
			result:  `{"valid":true}`,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ValidateSettings(ctx, tc.payload)
			require.NoError(t, err)
			assert.Equal(t, tc.result, string(result))
		})
	}
}
