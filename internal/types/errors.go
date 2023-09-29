package types

type MetricsHandlerInvalidType struct{}
type MetricsHandlerInvalidName struct{}
type MetricsHandlerInvalidValue struct{}
type UnknownScheme struct{}

func (e *MetricsHandlerInvalidType) Error() string {
	return "invalid metrics type"
}

func (e *MetricsHandlerInvalidName) Error() string {
	return "invalid metrics type"
}

func (e *MetricsHandlerInvalidValue) Error() string {
	return "invalid metrics type"
}

func (e *UnknownScheme) Error() string {
	return "invalid metrics type"
}

func NewMetricsHandlerInvalidTypeError() error {
	return &MetricsHandlerInvalidType{}
}

func NewMetricsHandlerInvalidNameError() error {
	return &MetricsHandlerInvalidName{}
}

func NewInvalidMetricsValue() error {
	return &MetricsHandlerInvalidValue{}
}

func NewUnknowSchemeError() error {
	return &UnknownScheme{}
}
