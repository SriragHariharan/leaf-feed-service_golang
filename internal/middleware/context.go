package middleware

import "context"

type contextKey int

const userIDKey contextKey = iota

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok || userID == "" {
		return "", false
	}
	return userID, true
}
