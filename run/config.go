package run

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/berquerant/fflist/query"
	"gopkg.in/yaml.v3"
)

func ParseConfig(r io.Reader) (*Config, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		if yErr := yaml.Unmarshal(b, &c); yErr != nil {
			return nil, errors.Join(err, yErr)
		}
	}

	if err := c.validate(); err != nil {
		return nil, err
	}
	return &c, nil
}

var (
	ErrConfig = errors.New("Config")
)

type Config struct {
	Root  []string   `json:"root" yaml:"root"`
	Query [][]string `json:"query" yaml:"query"`
}

func (c Config) validate() error {
	if len(c.Root) == 0 {
		return fmt.Errorf("%w: no root", ErrConfig)
	}
	if len(c.Query) == 0 {
		return fmt.Errorf("%w: no query", ErrConfig)
	}
	for i, x := range c.Query {
		if len(x) == 0 {
			return fmt.Errorf("%w: empty query at index %d", ErrConfig, i)
		}
	}
	return nil
}

func (c Config) ParseQuery() (query.Selector, error) {
	r := make([]query.Selector, len(c.Query))
	for i, a := range c.Query {
		s, err := ParseQuery(a)
		if err != nil {
			return nil, fmt.Errorf("%w: index %d", err, i)
		}
		r[i] = s
	}
	return query.NewOrSelector(r...), nil
}
