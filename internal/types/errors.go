package types

import "errors"

type MetricsHandlerInvalidType error
type MetricsHandlerInvalidName error
type MetricsHandlerInvalidValue error
type UnknownScheme error

func NewMetricsHandlerInvalidTypeError() MetricsHandlerInvalidType {
	return errors.New("invalid metrics type")
}

func NewMetricsHandlerInvalidNameError() MetricsHandlerInvalidName {
	return errors.New("invalid metrics name")
}

func NewInvalidMetricsValue() MetricsHandlerInvalidValue {
	return errors.New("invalid metrics value")
}

func NewUnknowSchemeError() error {
	return errors.New("unknown metrics scheme")
}
