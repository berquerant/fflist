package query

import (
	"context"
	"log/slog"
	"regexp"

	"github.com/berquerant/fflist/info"
	"github.com/berquerant/fflist/metric"
)

var (
	_ Selector = &RegexpSelector{}
)

type RegexpSelector struct {
	key string
	r   *regexp.Regexp
}

func NewRegexpSelector(q Query) (*RegexpSelector, error) {
	r, err := regexp.Compile(q.Value())
	if err != nil {
		return nil, err
	}
	return &RegexpSelector{
		key: q.Key(),
		r:   r,
	}, nil
}

func (s RegexpSelector) Select(_ context.Context, data info.Getter) bool {
	logAttr := []any{
		slog.String("key", s.key),
		slog.String("r", s.r.String()),
	}
	defer func() {
		slog.Debug("RegexpSelector", logAttr...)
	}()
	metric.IncrSelectCount()

	v, ok := data.Get(s.key)
	logAttr = append(logAttr, slog.Bool("found", ok))
	if !ok {
		logAttr = append(logAttr, slog.Bool("result", false))
		metric.IncrSelectDataMissingCount()
		metric.IncrSelectFailedCount()
		return false
	}
	r := s.r.MatchString(v)
	logAttr = append(logAttr, slog.String("value", v), slog.Bool("result", r))
	if r {
		metric.IncrSelectSuccessCount()
	} else {
		metric.IncrSelectFailedCount()
	}
	return r
}
