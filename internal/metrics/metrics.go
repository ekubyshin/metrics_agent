package metrics

import "database/sql"

type Gauge float64
type Counter int64

type Keyable[K any] interface {
	Key() K
	Serialize() map[string]any
}

type Metrics struct {
	ID    string   `json:"id" db:"id"`                 // имя метрики
	MType string   `json:"type" db:"type"`             // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty" db:"delta"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty" db:"value"` // значение метрики в случае передачи gauge
}

type MetricsKey struct {
	ID   string `json:"id" db:"id"`
	Type string `json:"type" db:"type"`
}

func (m Metrics) Key() MetricsKey {
	return MetricsKey{ID: m.ID, Type: m.MType}
}

func (m Metrics) Serialize() map[string]any {
	r := make(map[string]any)
	r["id"] = m.ID
	r["type"] = m.MType
	if m.Delta != nil {
		r["delta"] = *m.Delta
	} else {
		r["delta"] = sql.NullInt64{}
	}
	if m.Value != nil {
		r["value"] = *m.Value
	} else {
		r["value"] = sql.NullFloat64{}
	}
	return r
}
