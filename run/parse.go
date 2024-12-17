package run

import (
	"fmt"
	"os"

	"github.com/berquerant/fflist/query"
	"github.com/berquerant/fflist/slicesx"
)

const (
	queryShKey = "sh"
)

func ParseQuery(args []string) (query.Selector, error) {
	r := make([]query.Selector, len(args))
	for i, a := range args {
		a = os.ExpandEnv(a)
		x, err := query.Parse(a)
		if err != nil {
			return nil, fmt.Errorf("%w: index %d", err, i)
		}

		var s query.Selector
		switch x.Key() {
		case queryShKey:
			s = query.NewScriptSelector(x)
		default:
			s, err = query.NewRegexpSelector(x)
			if err != nil {
				return nil, fmt.Errorf("%w: index %d", err, i)
			}
		}
		r[i] = s
	}

	return query.NewAndSelector(r...), nil
}

func ParseQueryCommandLine(args []string) (query.Selector, error) {
	xs := slicesx.Chunk(args, "or", "OR")
	r := make([]query.Selector, len(xs))
	for i, x := range xs {
		s, err := ParseQuery(x)
		if err != nil {
			return nil, err
		}
		r[i] = s
	}

	return query.NewOrSelector(r...), nil
}

func ExpandEnvAll(ss ...string) []string {
	r := make([]string, len(ss))
	for i, s := range ss {
		r[i] = os.ExpandEnv(s)
	}
	return r
}
