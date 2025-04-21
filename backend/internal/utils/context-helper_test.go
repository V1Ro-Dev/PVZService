package utils_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"pvz/internal/utils"
	"pvz/pkg/logger"
)

func TestSetRequestId(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		expectErr bool
	}{
		{
			name:      "Test with empty context",
			ctx:       context.Background(),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxWithRequestID := utils.SetRequestId(tt.ctx)

			requestID, ok := ctxWithRequestID.Value(logger.RequestID).(string)

			assert.True(t, ok, "RequestId should be present in the context")
			assert.NotEmpty(t, requestID, "RequestId should not be empty")

			_, err := uuid.Parse(requestID)
			if tt.expectErr {
				assert.Error(t, err, "RequestId should be a valid UUID")
			} else {
				assert.NoError(t, err, "RequestId should be a valid UUID")
			}
		})
	}
}
