package container

import (
	"reflect"
	"sync"
)

type Container struct {
	mu    sync.RWMutex
	types map[reflect.Type]map[string]reflect.Value
}

func (c *Container) Register(name string, t reflect.Type, v any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.types == nil {
		c.types = make(map[reflect.Type]map[string]reflect.Value)
	}

	if _, ok := c.types[t]; !ok {
		c.types[t] = make(map[string]reflect.Value)
	}

	c.types[t][name] = reflect.ValueOf(v)
}

func (c *Container) Resolve(name string, t reflect.Type) (value any, found bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.types == nil {
		return nil, false
	}

	if _, ok := c.types[t]; !ok {
		return nil, false
	}

	v, ok := c.types[t][name]
	return v.Interface(), ok
}

func (c *Container) ResolveAll(t reflect.Type) (values []any) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.types == nil {
		return nil
	}

	if _, ok := c.types[t]; !ok {
		return nil
	}

	for _, v := range c.types[t] {
		values = append(values, v.Interface())
	}

	return values
}
