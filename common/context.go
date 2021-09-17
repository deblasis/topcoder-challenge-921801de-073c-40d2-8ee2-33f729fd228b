package common

import (
	"context"

	"github.com/google/uuid"
)

type contextKey int

const (
	ContextKeyUserId   = iota
	ContextKeyUserRole = iota

	ContextKeyReturnCode = iota

	ContextKeyRegisterShipOutcome contextKey = iota
	ContextKeyError                          = iota
)

func ExtractUserIdFromCtx(ctx context.Context) string {
	ctxUserId, ok := ctx.Value(ContextKeyUserId).(string)
	if !ok {
		return ""
	}
	userUuid, err := uuid.Parse(ctxUserId)
	if err != nil {
		return ""
	}
	return userUuid.String()
}

func ExtractUserRoleFromCtx(ctx context.Context) string {
	ctxUserRole, ok := ctx.Value(ContextKeyUserRole).(string)
	if !ok {
		return ""
	}
	return ctxUserRole
}
