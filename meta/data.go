package meta

import (
	"encoding/json"
	"maps"
)

type Data struct {
	d map[string]string
}

func NewData(d map[string]string) *Data {
	return &Data{
		d: d,
	}
}

func (d Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.d)
}

func (d Data) MarshalYAML() (any, error) {
	return d.d, nil
}

func (d Data) Get(key string) (string, bool) {
	x, ok := d.d[key]
	return x, ok
}

func (d Data) Merge(right *Data) *Data {
	if right == nil {
		return &d
	}
	left := maps.Clone(d.d)
	for k, v := range right.d {
		left[k] = v
	}
	return NewData(left)
}
