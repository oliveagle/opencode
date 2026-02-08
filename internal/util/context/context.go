package contextx

import (
	"context"
	"time"
)

// WithTimeout creates a context with timeout and returns a cancel function
func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

// WithCancel creates a cancellable context
func WithCancel(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(parent)
}

// WithDeadline creates a context with deadline
func WithDeadline(parent context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(parent, deadline)
}

// Background returns a background context
func Background() context.Context {
	return context.Background()
}

// TODO returns a TODO context
func TODO() context.Context {
	return context.TODO()
}

// Cause returns the cause of context cancellation
func Cause(ctx context.Context) error {
	return context.Cause(ctx)
}

// AfterFunc registers a function to call after context is done
func AfterFunc(ctx context.Context, f func()) context.CancelFunc {
	return context.AfterFunc(ctx, f)
}

// WithoutCancel returns a context that is not cancelled when parent is cancelled
func WithoutCancel(parent context.Context) context.Context {
	return context.WithoutCancel(parent)
}

// ValueCtx wraps a context with a value
type ValueCtx struct {
	context.Context
	key   interface{}
	value interface{}
}

// Value returns the value for key
func (c *ValueCtx) Value(key interface{}) interface{} {
	if key == c.key {
		return c.value
	}
	return c.Context.Value(key)
}

// WithValue returns a context with the given value
func WithValue(parent context.Context, key, value interface{}) context.Context {
	return &ValueCtx{Context: parent, key: key, value: value}
}

// StringKey is a typed string key for context values
type StringKey string

// Get retrieves a string value from context
func Get(ctx context.Context, key StringKey) (string, bool) {
	v := ctx.Value(key)
	if v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// GetInt retrieves an int value from context
func GetInt(ctx context.Context, key StringKey) (int, bool) {
	v := ctx.Value(key)
	if v == nil {
		return 0, false
	}
	i, ok := v.(int)
	return i, ok
}

// GetBool retrieves a bool value from context
func GetBool(ctx context.Context, key StringKey) (bool, bool) {
	v := ctx.Value(key)
	if v == nil {
		return false, false
	}
	b, ok := v.(bool)
	return b, ok
}

// Set sets a value in context and returns a new context
func Set(ctx context.Context, key StringKey, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}
