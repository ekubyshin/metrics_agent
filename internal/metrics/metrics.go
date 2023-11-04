package metrics

type Gauge float64
type Counter int64

type Keyable[K any] interface {
	Key() K
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricsKey struct {
	ID    string
	MType string
}

func (m Metrics) Key() MetricsKey {
	return MetricsKey{ID: m.ID, MType: m.MType}
}
