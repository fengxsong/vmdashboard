package request

import (
	"context"

	"github.com/fengxsong/vmdashboard/pkg/authentication/user"
)

// The key type is unexported to prevent collisions
type key int

const (
	// namespaceKey is the context key for the request namespace.
	namespaceKey key = iota

	// userKey is the context key for the request user.
	userKey
	errKey
)

// WithValue returns a copy of parent in which the value associated with key is val.
func WithValue(parent context.Context, key interface{}, val interface{}) context.Context {
	return context.WithValue(parent, key, val)
}

// WithNamespace returns a copy of parent in which the namespace value is set
func WithNamespace(parent context.Context, namespace string) context.Context {
	return WithValue(parent, namespaceKey, namespace)
}

// NamespaceFrom returns the value of the namespace key on the ctx
func NamespaceFrom(ctx context.Context) (string, bool) {
	namespace, ok := ctx.Value(namespaceKey).(string)
	return namespace, ok
}

// NamespaceValue returns the value of the namespace key on the ctx, or the empty string if none
func NamespaceValue(ctx context.Context) string {
	namespace, _ := NamespaceFrom(ctx)
	return namespace
}

// WithUser returns a copy of parent in which the user value is set
func WithUser(parent context.Context, user user.Info) context.Context {
	return WithValue(parent, userKey, user)
}

// UserFrom returns the value of the user key on the ctx
func UserFrom(ctx context.Context) (user.Info, bool) {
	user, ok := ctx.Value(userKey).(user.Info)
	return user, ok
}

func WithError(parent context.Context, err error) context.Context {
	return WithValue(parent, errKey, err)
}

func ErrorFrom(ctx context.Context) (error, bool) {
	err, ok := ctx.Value(errKey).(error)
	return err, ok
}
