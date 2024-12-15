package query_test

import (
	"context"
	"testing"

	"github.com/berquerant/fflist/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegexpSelector(t *testing.T) {
	t.Run("InvalidRegexp", func(t *testing.T) {
		_, err := query.NewRegexpSelector(query.NewQuery("key", "["))
		assert.NotNil(t, err)
	})

	for _, tc := range []struct {
		title string
		q     query.Query
		value string
		exist bool
		want  bool
	}{
		{
			title: "value matched",
			q:     query.NewQuery("k", "ref"),
			value: "ref",
			exist: true,
			want:  true,
		},
		{
			title: "value not matched",
			q:     query.NewQuery("k", "ref"),
			value: "fer",
			exist: true,
			want:  false,
		},
		{
			title: "value not found",
			q:     query.NewQuery("k", "ref"),
			exist: false,
			want:  false,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			s, err := query.NewRegexpSelector(tc.q)
			if !assert.Nil(t, err) {
				return
			}
			getter := new(mockInfoGetter)
			getter.On("Get", tc.q.Key()).Return(tc.value, tc.exist)
			got := s.Select(context.TODO(), getter)
			assert.Equal(t, tc.want, got)
		})
	}
}

type mockInfoGetter struct {
	mock.Mock
}

func (g *mockInfoGetter) Get(key string) (string, bool) {
	args := g.Called(key)
	return args.String(0), args.Bool(1)
}
