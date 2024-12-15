package info

import (
	"encoding/json"
	"fmt"
	"path/filepath"
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
	path := entry.Path()
	name := entry.Info().Name()
	return meta.NewData(map[string]string{
		"path":     path,
		"dir":      filepath.Dir(path),
		"name":     name,
		"ext":      filepath.Ext(name),
		"size":     fmt.Sprint(entry.Info().Size()),
		"mode":     fmt.Sprintf("%o", entry.Info().Mode()),
		"mod_time": entry.Info().ModTime().Format(time.DateTime),
	})
}
