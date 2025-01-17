package container

var globalContainer = &Container{}

func Register[T any](v T) {
	(&GenericContainer[T]{Container: globalContainer}).Register(v)
}

func RegisterNamed[T any](name string, v T) {
	(&GenericContainer[T]{Container: globalContainer}).RegisterNamed(name, v)
}

func Resolve[T any]() (T, bool) {
	return (&GenericContainer[T]{Container: globalContainer}).Resolve()
}

func ResolveNamed[T any](name string) (T, bool) {
	return (&GenericContainer[T]{Container: globalContainer}).ResolveNamed(name)
}

func ResolveAll[T any]() []T {
	return (&GenericContainer[T]{Container: globalContainer}).ResolveAll()
}
