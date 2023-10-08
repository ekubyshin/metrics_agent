package utils

func ToPointer[T any](val T) *T {
	return &val
}
