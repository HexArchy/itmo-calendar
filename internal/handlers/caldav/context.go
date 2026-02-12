package caldav

import "context"

type ctxKey struct{}

// ContextWithISU stores the ISU number in the context.
func ContextWithISU(ctx context.Context, isu int64) context.Context {
	return context.WithValue(ctx, ctxKey{}, isu)
}

// ISUFromContext retrieves the ISU number from the context.
func ISUFromContext(ctx context.Context) (int64, bool) {
	isu, ok := ctx.Value(ctxKey{}).(int64)
	return isu, ok
}
