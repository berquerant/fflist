package query_test

import (
	"context"
	"testing"

	"github.com/berquerant/fflist/info"
	"github.com/berquerant/fflist/meta"
	"github.com/berquerant/fflist/query"
	"github.com/stretchr/testify/assert"
)

func TestScriptSelector(t *testing.T) {
	for _, tc := range []struct {
		title string
		q     query.Query
		data  info.Getter
		want  bool
	}{
		{
			title: "grep hit",
			q:     query.NewQuery("sh", "grep -q ref"),
			data: info.New(meta.NewData(map[string]string{
				"value": "ref",
			})),
			want: true,
		},
		{
			title: "grep not hit",
			q:     query.NewQuery("sh", "grep -q ref"),
			data: info.New(meta.NewData(map[string]string{
				"value": "fer",
			})),
			want: false,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			s := query.NewScriptSelector(tc.q)
			got := s.Select(context.TODO(), tc.data)
			assert.Equal(t, tc.want, got)
		})
	}
}
