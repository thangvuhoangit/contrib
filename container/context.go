package container

import (
	"context"
	"errors"
)

var (
	ErrContainerNil = errors.New("Container not found")
)

type contextKey string

var containerContextKey = contextKey("ioc.container")

func WithContext(ctx context.Context, container *Container) context.Context {
	return context.WithValue(ctx, containerContextKey, container)
}

func FromContext(ctx context.Context) *Container {
	container, ok := ctx.Value(containerContextKey).(*Container)
	if !ok {
		panic(ErrContainerNil)
	}

	return container
}

func RegisterToCtx[T any](ctx context.Context, v T) {
	container := FromContext(ctx)
	(&GenericContainer[T]{Container: container}).Register(v)
}

func RegisterNamedToCtx[T any](ctx context.Context, name string, v T) {
	container := FromContext(ctx)
	(&GenericContainer[T]{Container: container}).RegisterNamed(name, v)
}

func ResolveFromCtx[T any](ctx context.Context) (T, bool) {
	container := FromContext(ctx)
	return (&GenericContainer[T]{Container: container}).Resolve()
}

func ResolveNamedFromCtx[T any](ctx context.Context, name string) (T, bool) {
	container := FromContext(ctx)
	return (&GenericContainer[T]{Container: container}).ResolveNamed(name)
}

func ResolveAllFromCtx[T any](ctx context.Context) []T {
	container := FromContext(ctx)
	return (&GenericContainer[T]{Container: container}).ResolveAll()
}
