package pointer

func From[T any](val T) *T {
	return &val
}
