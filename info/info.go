package info

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/berquerant/fflist/meta"
	"github.com/berquerant/fflist/walk"
)

type Getter interface {
	Get(key string) (string, bool)
}

var (
	_ Getter = &Metadata{}
)

func New(dataList ...*meta.Data) *Metadata {
	data := meta.NewData(map[string]string{})
	// flatten data
	for _, x := range dataList {
		data = data.Merge(x)
	}
	return &Metadata{
		data: data,
	}
}

type Metadata struct {
	data *meta.Data
}

func (d Metadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.data)
}

func (d Metadata) MarshalYAML() (any, error) {
	return d.data, nil
}

func (d Metadata) Get(key string) (string, bool) {
	return d.data.Get(key)
}

func NewMetadataFromEntry(entry walk.Entry) *meta.Data {
	return meta.NewData(map[string]string{
		"path":     entry.Path(),
		"name":     entry.Info().Name(),
		"size":     fmt.Sprint(entry.Info().Size()),
		"mode":     fmt.Sprintf("%o", entry.Info().Mode()),
		"mod_time": entry.Info().ModTime().Format(time.DateTime),
	})
}
