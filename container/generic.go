package container

import "reflect"

const _DEFAULT_ = ""

type GenericContainer[T any] struct {
	Container *Container
}

func (c *GenericContainer[T]) Register(v T) {
	c.RegisterNamed(_DEFAULT_, v)
}

func (c *GenericContainer[T]) RegisterNamed(name string, v T) {
	if c.Container == nil {
		panic("GenericContainer: Container is nil")
	}

	var instance T
	c.Container.Register(name, reflect.TypeOf(&instance).Elem(), v)
}

func (c *GenericContainer[T]) Resolve() (T, bool) {
	if c.Container == nil {
		panic("GenericContainer: Container is nil")
	}

	var instance T
	v, ok := c.Container.Resolve(_DEFAULT_, reflect.TypeOf(&instance).Elem())
	if !ok {
		return instance, false
	}

	return v.(T), ok
}

func (c *GenericContainer[T]) ResolveNamed(name string) (T, bool) {
	if c.Container == nil {
		panic("GenericContainer: Container is nil")
	}

	var instance T
	v, ok := c.Container.Resolve(name, reflect.TypeOf(&instance).Elem())
	if !ok {
		return instance, false
	}

	return v.(T), ok
}

func (c *GenericContainer[T]) ResolveAll() []T {
	if c.Container == nil {
		panic("GenericContainer: Container is nil")
	}

	var instance T
	values := c.Container.ResolveAll(reflect.TypeOf(&instance).Elem())
	if values == nil {
		return nil
	}

	var result []T
	for _, v := range values {
		result = append(result, v.(T))
	}

	return result
}
