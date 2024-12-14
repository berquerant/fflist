package query_test

import (
	"testing"

	"github.com/berquerant/fflist/query"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	for _, tc := range []struct {
		s    string
		want query.Query
		err  error
	}{
		{
			s:    "key=",
			want: query.NewQuery("key", ""),
		},
		{
			s:    "key=value",
			want: query.NewQuery("key", "value"),
		},
		{
			s:   "without_equal",
			err: query.ErrInvalidQuery,
		},
	} {
		t.Run(tc.s, func(t *testing.T) {
			got, err := query.Parse(tc.s)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
