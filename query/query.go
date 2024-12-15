package query

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/berquerant/fflist/info"
)

//go:generate go run github.com/berquerant/dataclass -type Query -field "Key string|Value string" -output query_dataclass_generated.go

var (
	ErrInvalidQuery = errors.New("InvalidQuery")
)

// Parse key=value string into Query.
func Parse(s string) (Query, error) {
	if !strings.Contains(s, "=") {
		return nil, fmt.Errorf("%w: %s", ErrInvalidQuery, s)
	}
	xs := strings.SplitN(s, "=", 2)
	return NewQuery(xs[0], xs[1]), nil
}

// Selector determines if a file matches the conditions using its metadata.
type Selector interface {
	Select(ctx context.Context, data info.Getter) bool
}

var (
	_ Selector = &AndSelector{}
)

func NewAndSelector(selectors ...Selector) *AndSelector {
	return &AndSelector{
		selectors: selectors,
	}
}

type AndSelector struct {
	selectors []Selector
}

func (s AndSelector) Select(ctx context.Context, data info.Getter) bool {
	for _, x := range s.selectors {
		if !x.Select(ctx, data) {
			return false
		}
	}
	return true
}

var (
	_ Selector = &OrSelector{}
)

func NewOrSelector(selectors ...Selector) *OrSelector {
	return &OrSelector{
		selectors: selectors,
	}
}

type OrSelector struct {
	selectors []Selector
}

func (s OrSelector) Select(ctx context.Context, data info.Getter) bool {
	for _, x := range s.selectors {
		if x.Select(ctx, data) {
			return true
		}
	}
	return false
}
