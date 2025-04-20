package utils

import (
	"context"

	"github.com/google/uuid"

	"pvz/pkg/logger"
)

func SetRequestId(ctx context.Context) context.Context {
	return context.WithValue(ctx,
		logger.RequestID,
		logger.ReqIdKey(uuid.New().String()))
}
