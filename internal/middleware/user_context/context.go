package user_context

import (
	"context"

	"github.com/google/uuid"
)

// UserContext encapsulates user-specific data.
type UserContext struct {
	UserUuid uuid.UUID
	// TODO fill
	UserName string
	Roles    []string
}

// contextKey is a type for context keys to avoid collisions.
type contextKey string

const userContextKey contextKey = "user_context"

// WithUserContext stores UserContext in the context.
func WithUserContext(ctx context.Context, uc UserContext) context.Context {
	return context.WithValue(ctx, userContextKey, uc)
}

// GetUserContext retrieves UserContext from the context.
func GetUserContext(ctx context.Context) (UserContext, bool) {
	uc, ok := ctx.Value(userContextKey).(UserContext)

	return uc, ok
}
